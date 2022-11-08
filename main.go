package main

import (
	"database/sql"
	"log"
	"net/http"

	"remi/internal/services"
	"remi/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Panicf("error loading config %v", err)
	}

	db, err := sql.Open(cfg.Postgres.ConnectionString())
	if err != nil {
		log.Panicf("error opening db %v", err)
	}

	remiService := services.NewRemiService(db, cfg.JWTSecret, cfg.URL)

	log.Printf("HTTP server listening at %v", cfg.HTTP.Address())
	err = http.ListenAndServe(cfg.HTTP.Address(), remiService)
	if err != nil {
		log.Panicf("error when starting server %v", err)
	}
}
