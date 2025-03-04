package main

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Message struct {
	Name string
	Body string
	Time int64
	Id   int64
}

type MessageParams struct {
	Name string
	Body string
}

type NotFoundError struct{}

type NotSavedError struct{}

type APIError struct{}

func (e *NotFoundError) Error() string {
	return "record not found"
}

func (e *NotSavedError) Error() string {
	return "record not saved"
}

func (e *APIError) Error() string {
	return "an error occurred"
}

func findAll(ctx context.Context) ([]Message, error) {
	rows, err := db.QueryContext(ctx, "SELECT * from messages")
	if err != nil {
		return nil, &APIError{}
	}
	var allMessages []Message
	for rows.Next() {
		var message Message
		rows.Scan(&message.Name, &message.Body, &message.Time, &message.Id)
		allMessages = append(allMessages, message)
	}
	defer rows.Close()
	return allMessages, nil
}

func findById(messageId string, ctx context.Context) (Message, error) {
	var message Message
	err := db.QueryRowContext(ctx, "SELECT * from messages where id = $1", messageId).Scan(&message.Name, &message.Body, &message.Time, &message.Id)

	if errors.Is(err, sql.ErrNoRows) {
		return Message{}, &NotFoundError{}
	} else if err != nil {
		return Message{}, &APIError{}
	}

	return message, nil
}

func updateMessage(message *Message, ctx context.Context) error {
	result, err := db.ExecContext(ctx, "UPDATE messages SET name = $1, body = $2 where id = $3", message.Name, message.Body, message.Id)
	if err != nil {
		return &APIError{}
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected < 1 {
		return &NotSavedError{}
	}
	return nil
}

func saveMessage(name string, body string, ctx context.Context) (Message, error) {
	message := Message{}
	ts := time.Now().Unix()
	row := db.QueryRowContext(ctx, "INSERT INTO messages values ($1, $2, $3) returning *", name, body, ts)

	err := row.Scan(&message.Name, &message.Body, &message.Time, &message.Id)
	if err != nil {
		return Message{}, &NotSavedError{}
	}
	return message, nil
}

func deleteMessage(messageId string, ctx context.Context) error {
	result, err := db.ExecContext(ctx, "DELETE FROM messages where id = $1", messageId)
	if err != nil {
		return &APIError{}
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected < 1 {
		return &NotFoundError{}
	}
	return nil
}
