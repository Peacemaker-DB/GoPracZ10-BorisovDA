package main

import (
	"log"
	"net/http"

	router "example.com/goprac10-borisovda/internal/http"
	"example.com/goprac10-borisovda/internal/platform/config"
)

func main() {
	cfg := config.Load()
	mux := router.Build(cfg)
	log.Println("listening on", cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, mux))
}
