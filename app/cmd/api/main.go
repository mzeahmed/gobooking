package main

import (
	"log"
	"net/http"

	"github.com/mzeahmed/go-booking/internal/config"
	"github.com/mzeahmed/go-booking/internal/router"
)

func main() {

	cfg := config.Load()

	r := router.New()

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	log.Printf("Server listening on :%s", cfg.Port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
