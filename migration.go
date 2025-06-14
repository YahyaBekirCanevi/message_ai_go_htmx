package main

import (
	"database/sql"
	"log"
)

// MigrateDB creates necessary tables in the SQLite database if they don't already exist.
func MigrateDB(db *sql.DB) error {
	// SQL statement to create the 'users' table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL
	);
	`

	log.Println("Attempting to create 'users' table...")
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Printf("Error creating users table: %v", err)
		return err
	}
	log.Println("'users' table created or already exists.")

	return nil
}
