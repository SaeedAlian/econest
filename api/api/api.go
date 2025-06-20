package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"github.com/swaggo/http-swagger"

	"github.com/SaeedAlian/econest/api/config"
	db_manager "github.com/SaeedAlian/econest/api/db/manager"
	_ "github.com/SaeedAlian/econest/api/docs"
	"github.com/SaeedAlian/econest/api/services/auth"
	"github.com/SaeedAlian/econest/api/services/product"
	"github.com/SaeedAlian/econest/api/services/smtp"
	"github.com/SaeedAlian/econest/api/services/store"
	"github.com/SaeedAlian/econest/api/services/user"
	"github.com/SaeedAlian/econest/api/services/wallet"
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

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(
			"http://localhost:5000/swagger/doc.json",
		),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods("GET")

	userSubrouter := router.PathPrefix("/user").Subrouter()
	storeSubrouter := router.PathPrefix("/store").Subrouter()
	productSubrouter := router.PathPrefix("/product").Subrouter()
	roleAndPermissionSubrouter := router.PathPrefix("/rp").Subrouter()
	walletSubrouter := router.PathPrefix("/wallet").Subrouter()
	orderSubrouter := router.PathPrefix("/order").Subrouter()

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

	storeService := store.NewHandler(dbManager, authHandler)
	storeService.RegisterRoutes(storeSubrouter)

	productService := product.NewHandler(dbManager, authHandler)
	productService.RegisterRoutes(productSubrouter)

	roleAndPermissionService := product.NewHandler(dbManager, authHandler)
	roleAndPermissionService.RegisterRoutes(roleAndPermissionSubrouter)

	walletService := wallet.NewHandler(dbManager, authHandler)
	walletService.RegisterRoutes(walletSubrouter)

	orderService := store.NewHandler(dbManager, authHandler)
	orderService.RegisterRoutes(orderSubrouter)

	log.Println("API Listening on ", s.addr)

	return http.ListenAndServe(s.addr, router)
}
