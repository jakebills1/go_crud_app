package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

var (
	db *sql.DB
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// db setup
	conninfo := os.Getenv("DATABASE_URL")
	connector, connErr := pq.NewConnector(conninfo)
	if connErr != nil {
		panic(connErr)
	}
	db = sql.OpenDB(connector)
	dbErr := db.Ping()
	if dbErr != nil {
		panic(dbErr)
	}
	defer db.Close()

	// http setup
	router := http.NewServeMux()
	router.HandleFunc("GET /messages/{messageId}", showHandler)
	router.HandleFunc("PUT /messages/{messageId}", updateHandler)
	router.HandleFunc("DELETE /messages/{messageId}", deleteHandler)
	router.HandleFunc("GET /messages/", indexHandler)
	router.HandleFunc("POST /messages/", createHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}
