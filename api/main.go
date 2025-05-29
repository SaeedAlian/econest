package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/SaeedAlian/econest/api/api"
	"github.com/SaeedAlian/econest/api/config"
	"github.com/SaeedAlian/econest/api/db"
	"github.com/SaeedAlian/econest/api/services/auth"
)

func main() {
	db, err := db.NewPGSQLStorage()
	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	ksCache := redis.NewClient(&redis.Options{
		Addr: config.Env.KeyServerRedisAddr,
	})

	keyServer := auth.NewKeyServer(ksCache)
	rotateKeys(keyServer)

	go func() {
		rotateHours := config.Env.RotateKeyDays * 24
		c := time.Tick(time.Duration(rotateHours) * time.Hour)
		for range c {
			rotateKeys(keyServer)
		}
	}()

	server := api.NewServer(fmt.Sprintf(":%s", config.Env.Port), db, keyServer)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Connection to DB was successful.")
}

func rotateKeys(keyServer *auth.KeyServer) {
	log.Println("rotating keys...")
	err := keyServer.RotateKeys(time.Now().String())
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("keys rotated")
}
