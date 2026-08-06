package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/akgbytes/ylx/internal/config"
	"github.com/akgbytes/ylx/internal/database"
	"github.com/akgbytes/ylx/internal/router"
)

const dbStartupTimeout = 30 * time.Second

func main() {

	cfg := config.MustLoad()

	dbCtx, dbCancel := context.WithTimeout(context.Background(), dbStartupTimeout)
	db, err := database.Connect(dbCtx, *cfg)
	dbCancel()
	if err != nil {
		log.Fatal("init database failed: ", err)
	}

	defer db.Close()
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

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
