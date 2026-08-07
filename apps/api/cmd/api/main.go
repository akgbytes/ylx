package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/akgbytes/ylx/internal/config"
	"github.com/akgbytes/ylx/internal/database"
	"github.com/akgbytes/ylx/internal/router"
)

const dbStartupTimeout = 30 * time.Second

func main() {
	if err := initServer(); err != nil {
		log.Fatal(err)
	}
}

func initServer() error {
	cfg := config.MustLoad()

	dbCtx, dbCancel := context.WithTimeout(context.Background(), dbStartupTimeout)
	db, err := database.Connect(dbCtx, *cfg)
	dbCancel()
	if err != nil {
		return fmt.Errorf("init database: %w", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("close db: %v", err)
		}
	}()
	log.Println("database connected")
	router := router.New()

	server := http.Server{
		Addr:         cfg.Addr,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.Println("starting ylx api server on port", server.Addr)

	return server.ListenAndServe()
}
