package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	id            int
	username      string
	name          string
	created_at    string
	password_hash string
}

type Message struct {
	id           int
	sender_id    int
	recipient_id int
	content      string
	sent_at      string
}

func getUser(db *sql.DB, userId int) (User, error) {
	row := db.QueryRow(
		"SELECT * FROM user WHERE id = ?",
		userId,
	)
	if row.Err() != nil {
		return User{}, row.Err()
	}

	user := User{}
	err := row.Scan(
		&user.id,
		&user.username,
		&user.name,
		&user.created_at,
		&user.password_hash,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func getMessages(db *sql.DB, userId1 int, userId2 int) ([]Message, error) {
	rows, err := db.Query(`
		SELECT * FROM message
		WHERE sender_id = ? AND recipient_id = ?
		OR sender_id = ? AND recipient_id = ?
		ORDER BY sent_at DESC
	`, userId1, userId2, userId2, userId1)
	if err != nil {
		return nil, err
	}

	messages := []Message{}
	for rows.Next() {
		var message Message
		err := rows.Scan(
			&message.id,
			&message.sender_id,
			&message.recipient_id,
			&message.content,
			&message.sent_at,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}
