package main

import (
	"log"
	"net/http"

	"github.com/akgbytes/ylx/internal/router"
)

func main() {

	router := router.New()

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("server starting on port 8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
