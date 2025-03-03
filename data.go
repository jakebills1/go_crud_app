package main

import (
	"database/sql"
	"errors"
	"log"
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

func (e *NotFoundError) Error() string {
	return "record not found"
}

func (e *NotSavedError) Error() string {
	return "record not saved"
}

func findAll() []Message {
	rows, _ := db.Query("SELECT * from messages")
	var allMessages []Message
	for rows.Next() {
		var message Message
		rows.Scan(&message.Name, &message.Body, &message.Time, &message.Id)
		allMessages = append(allMessages, message)
	}
	return allMessages
}

func findById(messageId string) (Message, error) {
	var message Message
	err := db.QueryRow("SELECT * from messages where id = $1", messageId).Scan(&message.Name, &message.Body, &message.Time, &message.Id)

	if errors.Is(err, sql.ErrNoRows) {
		return Message{}, &NotFoundError{}
	} else if err != nil {
		log.Fatal("QueryRow:", err)
	}

	return message, nil
}

func updateMessage(message *Message) error {
	result, err := db.Exec("UPDATE messages SET name = $1, body = $2 where id = $3", message.Name, message.Body, message.Id)
	if err != nil {
		log.Fatal("Exec():", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected < 1 {
		return &NotSavedError{}
	}
	return nil
}

func saveMessage(name string, body string) (Message, error) {
	ts := time.Now().Unix()
	row := db.QueryRow("INSERT INTO messages values ($1, $2, $3) returning id", name, body, ts)

	var id int64
	err := row.Scan(&id)
	log.Println(id)
	if err != nil {
		return Message{}, &NotSavedError{}
	}
	log.Println(id)
	message := Message{Name: name, Body: body, Time: ts, Id: id}
	return message, nil
}

func deleteMessage(messageId string) error {
	result, execErr := db.Exec("DELETE FROM messages where id = $1", messageId)
	if execErr != nil {
		log.Fatal("Exec():", execErr)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected < 1 {
		return &NotFoundError{}
	}
	return nil
}
