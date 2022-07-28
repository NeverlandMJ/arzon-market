package main

import (
	"fmt"
	"net/http"
	"os"

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

	http.ListenAndServe(GetPort(), r)
}

// Get the Port from the environment so we can run on Heroku
func GetPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "1323"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}