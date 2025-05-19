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
	subrouter := router.PathPrefix("/api").Subrouter()

	userSubrouter := subrouter.PathPrefix("/user").Subrouter()

	authCache := redis.NewClient(&redis.Options{
		Addr: config.Env.KeyServerRedisAddr,
	})

	authHandler := auth.NewAuthHandler(authCache, s.keyServer)
	dbManager := db_manager.NewManager(s.db)

	userService := user.NewHandler(dbManager, authHandler)
	userService.RegisterRoutes(userSubrouter)

	log.Println("API Listening on ", s.addr)

	return http.ListenAndServe(s.addr, router)
}
