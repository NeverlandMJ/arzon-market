package main

import (
	"arzon-market/postgres"
	"arzon-market/server"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	db, err := postgres.Connect()
	if err != nil {
		panic(err)
	}

	repo := postgres.NewPostgresRepository(db)
	r := server.NewRouter(repo)

	http.ListenAndServe(GetPort(), r)
}

// Get the Port from the environment so we can run on Heroku
func GetPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "8080"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}