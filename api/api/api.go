package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"

	"github.com/SaeedAlian/econest/api/config"
	db_manager "github.com/SaeedAlian/econest/api/db/manager"
	"github.com/SaeedAlian/econest/api/services/auth"
	"github.com/SaeedAlian/econest/api/services/product"
	"github.com/SaeedAlian/econest/api/services/smtp"
	"github.com/SaeedAlian/econest/api/services/store"
	"github.com/SaeedAlian/econest/api/services/user"
)

type Server struct {
	addr      string
	db        *sql.DB
	keyServer *auth.KeyServer
}

func NewServer(addr string, db *sql.DB, keyServer *auth.KeyServer) *Server {
	return &Server{
		addr:      addr,
		db:        db,
		keyServer: keyServer,
	}
}

func (s *Server) Run() error {
	router := mux.NewRouter()

	userSubrouter := router.PathPrefix("/user").Subrouter()
	storeSubrouter := router.PathPrefix("/store").Subrouter()
	productSubrouter := router.PathPrefix("/product").Subrouter()
	roleAndPermissionSubrouter := router.PathPrefix("/rp").Subrouter()

	authCache := redis.NewClient(&redis.Options{
		Addr: config.Env.KeyServerRedisAddr,
	})

	authHandler := auth.NewAuthHandler(authCache, s.keyServer)
	dbManager := db_manager.NewManager(s.db)
	smtpServer := smtp.NewSMTPServer(
		config.Env.SMTPHost,
		config.Env.SMTPPort,
		config.Env.SMTPEmail,
		config.Env.SMTPPassword,
	)

	userService := user.NewHandler(dbManager, authHandler, smtpServer)
	userService.RegisterRoutes(userSubrouter)

	storeService := store.NewHandler(dbManager, authHandler, smtpServer)
	storeService.RegisterRoutes(storeSubrouter)

	productService := product.NewHandler(dbManager, authHandler, smtpServer)
	productService.RegisterRoutes(productSubrouter)

	roleAndPermissionService := product.NewHandler(dbManager, authHandler, smtpServer)
	roleAndPermissionService.RegisterRoutes(roleAndPermissionSubrouter)

	log.Println("API Listening on ", s.addr)

	return http.ListenAndServe(s.addr, router)
}
