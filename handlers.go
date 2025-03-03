package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func indexHandler(w http.ResponseWriter, req *http.Request) {
	allMessages := findAll()
	b, err := json.Marshal(allMessages)
	if err != nil {
		log.Println("Marshal():", err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func createHandler(w http.ResponseWriter, req *http.Request) {
	body, readErr := io.ReadAll(req.Body)
	if readErr != nil {
		log.Fatal("ReadAll():", readErr)
	}
	var messageParams MessageParams
	unMarshalErr := json.Unmarshal(body, &messageParams)
	if unMarshalErr != nil {
		log.Fatal("Unmarshal():", unMarshalErr)
	}
	message, saveErr := saveMessage(messageParams.Name, messageParams.Body)
	if saveErr != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusCreated)
		b, marshalErr := json.Marshal(message)
		if marshalErr != nil {
			log.Fatal("Marshal():", marshalErr)
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(b)
	}
}

func showHandler(w http.ResponseWriter, req *http.Request) {
	requestedMessageId := req.PathValue("messageId")
	foundMessage, err := findById(requestedMessageId)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Header().Add("Content-Type", "application/json")
		b, err := json.Marshal(foundMessage)
		if err != nil {
			log.Fatal("Marshal():", err)
		}
		w.Write(b)
	}
}

func updateHandler(w http.ResponseWriter, req *http.Request) {
	requestedMessageId := req.PathValue("messageId")
	message, err := findById(requestedMessageId)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		body, readErr := io.ReadAll(req.Body)
		if readErr != nil {
			log.Println("ReadAll():", readErr)
		}
		unMarshalErr := json.Unmarshal(body, &message)
		if unMarshalErr != nil {
			log.Println("Unmarshal():", unMarshalErr)
		}
		updateErr := updateMessage(&message)
		if updateErr != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.Header().Add("Content-Type", "application/json")
			b, marshalErr := json.Marshal(message)
			if marshalErr != nil {
				log.Println("Marshal():", marshalErr)
			}
			w.Write(b)
		}
	}
}

func deleteHandler(w http.ResponseWriter, req *http.Request) {
	requestedMessageId := req.PathValue("messageId")
	err := deleteMessage(requestedMessageId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
