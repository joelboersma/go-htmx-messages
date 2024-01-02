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

func (u *User) getConversationList(db *sql.DB) ([]ConversationSummary, error) {
	query := `
		SELECT
			CASE
				WHEN message.sender_id = ? THEN message.recipient_id
				WHEN message.recipient_id = ? THEN message.sender_id
			END AS other_user_id,
			user.name AS other_user_name,
			MAX(message.sent_at) AS latest_interaction_time,
			message.content,
			CASE 
				WHEN message.recipient_id = 1 THEN TRUE
				ELSE FALSE
			END AS message_from_other_user
		FROM message
		JOIN user ON other_user_id = user.id
		WHERE message.sender_id = ? OR message.recipient_id = ?
		GROUP BY other_user_id
		ORDER BY latest_interaction_time DESC
	`
	rows, err := db.Query(query, u.id, u.id, u.id, u.id)
	if err != nil {
		return []ConversationSummary{}, err
	}
	defer rows.Close()

	conversationSummaries := []ConversationSummary{}
	for rows.Next() {
		conversationSummary := ConversationSummary{
			// Attrbitues from `u`
			u.id, u.name,
			// Other user and message info to be filled in from row
			0, "", "", "", false,
		}
		err = rows.Scan(
			&conversationSummary.OtherUserId,
			&conversationSummary.OtherUserName,
			&conversationSummary.LastMessageTimestamp,
			&conversationSummary.LastMessageContent,
			&conversationSummary.LastMessageFromOtherUser,
		)
		if err != nil {
			return []ConversationSummary{}, err
		}
		conversationSummaries = append(conversationSummaries, conversationSummary)
	}
	if err := rows.Err(); err != nil {
		return []ConversationSummary{}, err
	}
	return conversationSummaries, nil
}

type Message struct {
	id           int
	sender_id    int
	recipient_id int
	content      string
	sent_at      string
}

type ConversationSummary struct {
	UserId                   int
	UserName                 string
	OtherUserId              int
	OtherUserName            string
	LastMessageTimestamp     string
	LastMessageContent       string
	LastMessageFromOtherUser bool
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
		ORDER BY sent_at
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
