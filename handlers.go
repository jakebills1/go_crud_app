package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

func indexHandler(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()
	allMessages, err := findAll(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	b := marshal(allMessages)
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func createHandler(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()
	body := getBody(req)
	var messageParams Message
	parseBodyAsJson(body, &messageParams)
	message, saveErr := saveMessage(messageParams.Name, messageParams.Body, ctx)
	if saveErr != nil {
		http.Error(w, saveErr.Error(), http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusCreated)
		b := marshal(message)
		w.Header().Add("Content-Type", "application/json")
		w.Write(b)
	}
}

func showHandler(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()

	requestedMessageId := req.PathValue("messageId")
	foundMessage, err := findById(requestedMessageId, ctx)

	switch err.(type) {
	case nil:
		w.Header().Add("Content-Type", "application/json")
		b := marshal(foundMessage)
		w.Write(b)
	case *NotFoundError:
		http.Error(w, err.Error(), http.StatusNotFound)
	case *APIError:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func updateHandler(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()

	requestedMessageId := req.PathValue("messageId")
	message, err := findById(requestedMessageId, ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		body := getBody(req)
		parseBodyAsJson(body, &message)
		updateErr := updateMessage(&message, ctx)
		if updateErr != nil {
			http.Error(w, updateErr.Error(), http.StatusBadRequest)
		} else {
			w.Header().Add("Content-Type", "application/json")
			b := marshal(message)
			w.Write(b)
		}
	}
}

func deleteHandler(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()

	requestedMessageId := req.PathValue("messageId")
	err := deleteMessage(requestedMessageId, ctx)
	switch err.(type) {
	case nil:
		w.WriteHeader(http.StatusNoContent)
	case *NotFoundError:
		http.Error(w, err.Error(), http.StatusNotFound)
	case *APIError:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getBody(req *http.Request) []byte {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatal("ReadAll():", err)
	}
	return body
}

func parseBodyAsJson(body []byte, src *Message) {
	err := json.Unmarshal(body, src)
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
