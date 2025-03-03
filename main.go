package main

import (
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

var (
	db *sql.DB
)

func main() {
	setUpEnv()
	setUpDb()
	defer db.Close()

	log.Fatal(http.ListenAndServe(":8080", configureRouter()))
}

func setUpEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func setUpDb() {
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
}

func configureRouter() *http.ServeMux {
	// http setup
	router := http.NewServeMux()
	router.HandleFunc("GET /messages/{messageId}", showHandler)
	router.HandleFunc("PUT /messages/{messageId}", updateHandler)
	router.HandleFunc("DELETE /messages/{messageId}", deleteHandler)
	router.HandleFunc("GET /messages/", indexHandler)
	router.HandleFunc("POST /messages/", createHandler)
	return router
}
