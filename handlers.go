package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func indexHandler(w http.ResponseWriter, req *http.Request) {
	allMessages := findAll()
	b := marshal(allMessages)
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func createHandler(w http.ResponseWriter, req *http.Request) {
	body := getBody(req)
	var messageParams Message
	parseBodyAsJson(body, messageParams)
	message, saveErr := saveMessage(messageParams.Name, messageParams.Body)
	if saveErr != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusCreated)
		b := marshal(message)
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
		b := marshal(foundMessage)
		w.Write(b)
	}
}

func updateHandler(w http.ResponseWriter, req *http.Request) {
	requestedMessageId := req.PathValue("messageId")
	message, err := findById(requestedMessageId)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		body := getBody(req)
		parseBodyAsJson(body, message)
		updateErr := updateMessage(&message)
		if updateErr != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.Header().Add("Content-Type", "application/json")
			b := marshal(message)
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

func getBody(req *http.Request) []byte {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatal("ReadAll():", err)
	}
	return body
}

func parseBodyAsJson(body []byte, src Message) {
	err := json.Unmarshal(body, &src)
	if err != nil {
		log.Fatal("Unmarshal():", err)
	}
}

func marshal[T Message | []Message](src T) []byte {
	b, err := json.Marshal(src)
	if err != nil {
		log.Fatal("Marshal():", err)
	}
	return b
}
