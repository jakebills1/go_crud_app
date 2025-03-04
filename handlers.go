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
	b, marshalErr := marshal(allMessages)
	if marshalErr != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func createHandler(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()
	body, readErr := getBody(req)
	if readErr != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
	}
	var messageParams Message
	parseErr := parseBodyAsJson(body, &messageParams)
	if parseErr != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
	}
	message, saveErr := saveMessage(ctx, messageParams.Name, messageParams.Body)
	if saveErr != nil {
		http.Error(w, saveErr.Error(), http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusCreated)
		b, marshalErr := marshal(message)
		if marshalErr != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(b)
	}
}

func showHandler(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()

	requestedMessageId := req.PathValue("messageId")
	foundMessage, err := findById(ctx, requestedMessageId)

	switch err.(type) {
	case nil:
		w.Header().Add("Content-Type", "application/json")
		b, marshalErr := marshal(foundMessage)
		if marshalErr != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
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
	message, err := findById(ctx, requestedMessageId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		body, readErr := getBody(req)
		if readErr != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
		}
		parseErr := parseBodyAsJson(body, &message)
		if parseErr != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
		}
		updateErr := updateMessage(ctx, &message)
		if updateErr != nil {
			http.Error(w, updateErr.Error(), http.StatusBadRequest)
		} else {
			w.Header().Add("Content-Type", "application/json")
			b, marshalErr := marshal(message)
			if marshalErr != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			w.Write(b)
		}
	}
}

func deleteHandler(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()

	requestedMessageId := req.PathValue("messageId")
	err := deleteMessage(ctx, requestedMessageId)
	switch err.(type) {
	case nil:
		w.WriteHeader(http.StatusNoContent)
	case *NotFoundError:
		http.Error(w, err.Error(), http.StatusNotFound)
	case *APIError:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getBody(req *http.Request) ([]byte, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println("ReadAll():", err)
		return nil, err
	}
	return body, nil
}

func parseBodyAsJson(body []byte, src *Message) error {
	err := json.Unmarshal(body, src)
	if err != nil {
		log.Println("Unmarshal():", err)
		return err
	}
	return nil
}

func marshal[T Message | []Message](src T) ([]byte, error) {
	b, err := json.Marshal(src)
	if err != nil {
		log.Println("Marshal():", err)
		return nil, err
	}
	return b, nil
}
