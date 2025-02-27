package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Message struct {
	Name string
	Body string
	Time int64
}

var (
	m = Message{"Alice", "Hello", 1294706395881547000}
)

func main() {
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			var b, err = json.Marshal(m)
			if err != nil {
				log.Println("Marshal():", err)
			}
			w.Header().Add("Content-Type", "application/json")
			w.Write(b)
		} else if req.Method == "POST" {
			body, err := io.ReadAll(req.Body)
			if err != nil {
				log.Println("ReadAll():", err)
			}
			var m Message
			err = json.Unmarshal(body, &m)
			if err != nil {
				log.Println("Unmarshal():", err)
			}
			log.Println("creating new Message with Name =", m.Name, ", Body =", m.Body, ", and Time =", m.Time)
			// todo create record
			w.Header().Add("Content-Type", "text/html")
			w.WriteHeader(http.StatusCreated)
			responseBody := []byte("<html><body>Created new message!</body></html>")
			w.Write(responseBody)
		}
	}

	http.HandleFunc("/hello", helloHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
