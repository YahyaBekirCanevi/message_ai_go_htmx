package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Conversation struct {
	ID        int
	UserID    *int  
	CreatedAt time.Time
	Title     string
}

func FindConversationByID(db *sql.DB, id int) (*Conversation, error) {
    row := db.QueryRow("SELECT id, title FROM conversations WHERE id = ?", id)
    var conv Conversation
    err := row.Scan(&conv.ID, &conv.Title)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("conversation not found")
        }
        return nil, err
    }
    return &conv, nil
}

func (c *Conversation) CreateChat(db *sql.DB, firstTopic string) (int64, error) {
	c.Title = firstTopic
	c.CreatedAt = time.Now()

	stmt, err := db.Prepare("INSERT INTO conversations(user_id, created_at, title) VALUES(?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("error preparing insert statement: %w", err)
	}
	defer stmt.Close()

	var userID interface{}
	if c.UserID != nil {
		userID = *c.UserID
	} else {
		userID = nil
	}

	res, err := stmt.Exec(userID, c.CreatedAt, c.Title)
	if err != nil {
		return 0, fmt.Errorf("error executing insert statement: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert ID: %w", err)
	}
	c.ID = int(id)
	log.Printf("Conversation created with ID: %d, Title: %s", c.ID, c.Title)
	return id, nil
}

func GetAllChats(db *sql.DB) ([]Conversation, error) {
	rows, err := db.Query("SELECT id, user_id, created_at, title FROM conversations ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("error querying all conversations: %w", err)
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		var conv Conversation
		var userID sql.NullInt64
		var createdAt string
		if err := rows.Scan(&conv.ID, &userID, &createdAt, &conv.Title); err != nil {
			log.Printf("Error scanning conversation row: %v", err)
			continue
		}
		if userID.Valid {
			uid := int(userID.Int64)
			conv.UserID = &uid
		} else {
			conv.UserID = nil
		}
		conv.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		conversations = append(conversations, conv)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over conversation rows: %w", err)
	}
	return conversations, nil
}

func DeleteChat(db *sql.DB, id int) (int64, error) {
	stmt, err := db.Prepare("DELETE FROM conversations WHERE id = ?")
	if err != nil {
		return 0, fmt.Errorf("error preparing delete statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return 0, fmt.Errorf("error executing delete statement: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error getting rows affected: %w", err)
	}
	log.Printf("Conversation with ID %d deleted. Rows affected: %d", id, rowsAffected)
	return rowsAffected, nil
}
