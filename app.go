package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"time"
)

type Message struct {
	Name string
	Body string
	Time int64
	Id   uuid.UUID
}

var (
	messages = make([]Message, 0)
)

func main() {
	indexHandler := func(w http.ResponseWriter, req *http.Request) {
		b, err := json.Marshal(messages)
		if err != nil {
			log.Println("Marshal():", err)
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(b)
	}

	createHandler := func(w http.ResponseWriter, req *http.Request) {
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
		if m.Time == 0 {
			m.Time = time.Now().Unix()
		}
		m.Id = uuid.New()
		messages = append(messages, m)
		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(http.StatusCreated)
		responseBody := []byte("<html><body>Created new message!</body></html>")
		w.Write(responseBody)

	}

	showHandler := func(w http.ResponseWriter, req *http.Request) {
		requestedMessageId := req.PathValue("messageId")
		var foundMessage Message
		for _, message := range messages {
			if message.Id.String() == requestedMessageId {
				foundMessage = message
			}
		}

		if foundMessage.Time != 0 { // todo there is a better way to check if a message was found
			w.Header().Add("Content-Type", "application/json")
			b, err := json.Marshal(foundMessage)
			if err != nil {
				log.Println("Marshal():", err)
			}
			w.Write(b)
		} else {
			w.Header().Add("Content-Type", "text/html")
			w.WriteHeader(http.StatusNotFound)
			responseBody := []byte("<html><body>Message Not Found!</body></html>")
			w.Write(responseBody)
		}
	}

	updateHandler := func(w http.ResponseWriter, req *http.Request) {
		requestedMessageId := req.PathValue("messageId")
		var foundMessage Message
		for _, message := range messages {
			if message.Id.String() == requestedMessageId {
				foundMessage = message
			}
		}

		if foundMessage.Time != 0 {
			body, err := io.ReadAll(req.Body)
			if err != nil {
				log.Println("ReadAll():", err)
			}
			var updates Message
			err = json.Unmarshal(body, &updates)
			if err != nil {
				log.Println("Unmarshal():", err)
			}
			updatedMessages := make([]Message, len(messages))
			for _, message := range messages {
				if message.Id.String() != requestedMessageId {
					updatedMessages = append(updatedMessages, message)
				} else {
					if updates.Name != "" {
						message.Name = updates.Name
					}
					if updates.Body != "" {
						message.Body = updates.Body
					}
					updatedMessages = append(updatedMessages, message)
				}
			}
			messages = updatedMessages
			w.Header().Add("Content-Type", "text/html")
			responseBody := []byte("<html><body>Message updated/</body></html>")
			w.Write(responseBody)
		} else {
			w.Header().Add("Content-Type", "text/html")
			w.WriteHeader(http.StatusNotFound)
			responseBody := []byte("<html><body>Message Not Found!</body></html>")
			w.Write(responseBody)
		}

	}

	deleteHandler := func(w http.ResponseWriter, req *http.Request) {
		requestedMessageId := req.PathValue("messageId")
		var foundMessage Message
		for _, message := range messages {
			if message.Id.String() == requestedMessageId {
				foundMessage = message
			}
		}

		if foundMessage.Time != 0 {
			filteredMessages := make([]Message, len(messages))
			for _, message := range messages {
				if message.Id.String() != requestedMessageId {
					filteredMessages = append(filteredMessages, message)
				}
			}
			messages = filteredMessages
			w.Header().Add("Content-Type", "text/html")
			responseBody := []byte("<html><body>Message deleted/</body></html>")
			w.Write(responseBody)
		} else {
			w.Header().Add("Content-Type", "text/html")
			w.WriteHeader(http.StatusNotFound)
			responseBody := []byte("<html><body>Message Not Found!</body></html>")
			w.Write(responseBody)
		}

	}

	router := http.NewServeMux()
	router.HandleFunc("GET /messages/{messageId}", showHandler)
	router.HandleFunc("PUT /messages/{messageId}", updateHandler)
	router.HandleFunc("DELETE /messages/{messageId}", deleteHandler)
	router.HandleFunc("GET /messages/", indexHandler)
	router.HandleFunc("POST /messages/", createHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}
