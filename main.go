package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lib/pq"
)

var (
	db *sql.DB
)

func main() {
	// db setup
	conninfo := fmt.Sprintf("postgresql://%s@localhost:5432/go_crud_app_development?sslmode=disable", os.Getenv("USERNAME"))
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
