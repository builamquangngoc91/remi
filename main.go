package main

import (
	"database/sql"
	"log"
	"net/http"

	"remi/internal/services"
	"remi/pkg/config"

	"github.com/pressly/goose/v3"
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

	if err := goose.SetDialect("postgres"); err != nil {
		log.Panicf("goose.SetDialect: %w", err)
	}

	if err := goose.Up(db, "migrations/sql"); err != nil {
		log.Panicf("goose.Up: %w", err)
	}

	remiService := services.NewRemiService(db, cfg.JWTSecret, cfg.URL)

	log.Printf("HTTP server listening at %v", cfg.HTTP.Address())

	if err := http.ListenAndServe(cfg.HTTP.Address(), remiService); err != nil {
		log.Panicf("error when starting server %v", err)
	}
}
