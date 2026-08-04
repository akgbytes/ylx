package main

import (
	"log"
	"net/http"

	"github.com/akgbytes/ylx/internal/config"
	"github.com/akgbytes/ylx/internal/router"
)

func main() {

	cfg := config.MustLoad()

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
