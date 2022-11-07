package main

import (
	"database/sql"
	"log"
	"net/http"

	"remi/config"
	"remi/internal/services"
	cc "remi/pkg/config"
)

func main() {
	cc.InitFlags()
	cc.ParseFlags()

	cfg, err := config.Load()
	if err != nil {
		log.Panicf("error loading config %v", err)
	}

	db, err := sql.Open(cfg.Postgres.ConnectionString())
	if err != nil {
		log.Panicf("error opening db %v", err)
	}

	remiService := services.NewRemiService(db, cfg.JWTSecret)

	log.Printf("HTTP server listening at %v", cfg.HTTP.Address())
	err = http.ListenAndServe(cfg.HTTP.Address(), remiService)
	if err != nil {
		log.Panicf("error when starting server %v", err)
	}
}
