package main

import (
	"database/sql"
	"log"
)

func MigrateDB(db *sql.DB) error {
	if err := CreateConversationsTable(db); err != nil {
		return err
	}
	if err := CreateMessagesTable(db); err != nil {
		return err
	}
	return nil
}

func CreateConversationsTable(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS conversations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER, -- Optional: nullable user id for future user tracking
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		title TEXT NOT NULL UNIQUE
	);
	`

	log.Println("Attempting to create 'conversations' table...")
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Printf("Error creating conversations table: %v", err)
		return err
	}
	log.Println("'conversations' table created or already exists.")
	return nil
}

func CreateMessagesTable(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		conversation_id INTEGER NOT NULL,
		sender TEXT NOT NULL, -- "user" or "ai"
		content TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (conversation_id) REFERENCES conversations(id)
	);
	`

	log.Println("Attempting to create 'messages' table...")
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Printf("Error creating messages table: %v", err)
		return err
	}
	log.Println("'messages' table created or already exists.")
	return nil
}
