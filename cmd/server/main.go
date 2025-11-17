package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"hearts/internal/api"
	"hearts/internal/db"
)

func main() {
	dbConn, err := db.Open(os.Getenv("DB_CONN"))
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()
	err = db.CreateSchema(dbConn)
	if err != nil {
		panic(err)
	}

	// TODO: health endpoint should also test db connectivity
	http.Handle("/health",
		api.CommonHeaders(
			http.HandlerFunc(api.HealthHandler)))

	http.Handle("/match/{id}",
		api.WithDatabase(
			dbConn,
			api.CommonHeaders(
				http.HandlerFunc(api.MatchStateHandler))))

	http.Handle("/create-match",
		api.WithDatabase(
			dbConn,
			api.CommonHeaders(
				http.HandlerFunc(api.CreateMatchHandler))))

	http.Handle(
		"/join-match",
		api.WithDatabase(
			dbConn,
			api.CommonHeaders(http.HandlerFunc(api.JoinMatchHandler))))

	fmt.Println("listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
