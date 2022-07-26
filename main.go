package main

import (
	"net/http"
	"arzon-market/postgres"
	"arzon-market/server"

	_ "github.com/lib/pq"
)

func main() {
	db, err := postgres.Connect()
	if err != nil {
		panic(err)
	}

	repo := postgres.NewPostgresRepository(db)
	r := server.NewRouter(repo)

	http.ListenAndServe(":8080", r)
}