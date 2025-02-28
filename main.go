package main

import (
	"github.com/google/uuid"
	"log"
	"net/http"
)

type Message struct {
	Name string
	Body string
	Time int64
	Id   uuid.UUID
}

// todo use DB
var (
	messages = make([]Message, 0)
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("GET /messages/{messageId}", showHandler)
	router.HandleFunc("PUT /messages/{messageId}", updateHandler)
	router.HandleFunc("DELETE /messages/{messageId}", deleteHandler)
	router.HandleFunc("GET /messages/", indexHandler)
	router.HandleFunc("POST /messages/", createHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}
