package main

import (
	"github.com/NeverlandMJ/arzon-market/api"
	"github.com/NeverlandMJ/arzon-market/config"
	"github.com/NeverlandMJ/arzon-market/postgres"
	"github.com/NeverlandMJ/arzon-market/service"
	"github.com/NeverlandMJ/arzon-market/storage"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	db, err := postgres.Connect(cfg)
	if err != nil {
		panic(err)
	}

	repo := storage.NewPostgresRepository(db)
	serve := service.NewService(repo)
	r := api.NewRouter(serve)

	r.Run()
}