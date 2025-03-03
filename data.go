package main

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	Name string
	Body string
	Time int64
	Id   uuid.UUID
}

type MessageUpdates struct {
	Name string
	Body string
}

type NotFoundError struct{}

func (e *NotFoundError) Error() string {
	return "record not found"
}

func findById(messageId string) (Message, error) {
	for _, message := range messages {
		if message.Id.String() == messageId {
			return message, nil
		}
	}
	return Message{}, &NotFoundError{}
}

func updateMessage(message Message, updates MessageUpdates) Message {
	if updates.Name != "" {
		message.Name = updates.Name
	}
	if updates.Body != "" {
		message.Body = updates.Body
	}
	return message
}

func saveMessage(message *Message) {
	if message.Time == 0 {
		message.Time = time.Now().Unix()
	}
	message.Id = uuid.New()
}
