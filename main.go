package main

import (
	"github.com/NeverlandMJ/arzon-market/config"
	"github.com/NeverlandMJ/arzon-market/postgres"
	"github.com/NeverlandMJ/arzon-market/server"
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

	repo := server.NewPostgresRepository(db)
	r := server.NewRouter(repo)

	r.Run()
}
