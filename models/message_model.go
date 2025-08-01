package models

import (
	"database/sql"
	"time"
)

type Message struct {
	ID             int       `db:"id" json:"id"`
	ConversationID int       `db:"conversation_id" json:"conversation_id"`
	Sender         string    `db:"sender" json:"sender"` // "user" or "ai"
	Content        string    `db:"content" json:"content"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

func CreateMessage(db *sql.DB, conversationID int, sender string, content string) (error) {
	now := time.Now()
	_, err := db.Exec(
		`INSERT INTO messages (conversation_id, sender, content, created_at) VALUES (?, ?, ?, ?)`,
		conversationID, sender, content, now,
	)
	if err != nil {
		return err
	}
	return  nil
}

func GetMessagesByConversation(db *sql.DB, conversationID int) ([]Message, error) {
	rows, err := db.Query(
		`SELECT id, conversation_id, sender, content, created_at FROM messages WHERE conversation_id = ? ORDER BY created_at ASC`,
		conversationID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.Sender, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}