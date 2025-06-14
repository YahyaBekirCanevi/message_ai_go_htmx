package models

import (
	"database/sql"
	"fmt"
	"log"
)

type User struct {
	ID    int
	Name  string
	Email string
}

func (u *User) Insert(db *sql.DB) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO users(name, email) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("error preparing insert statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(u.Name, u.Email)
	if err != nil {
		return 0, fmt.Errorf("error executing insert statement: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert ID: %w", err)
	}
	u.ID = int(id) // Update the User struct with the new ID
	log.Printf("User '%s' (Email: %s) inserted with ID: %d", u.Name, u.Email, u.ID)
	return id, nil
}

func FindByID(db *sql.DB, id int) (*User, error) {
	row := db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %d not found", id)
		}
		return nil, fmt.Errorf("error scanning user with ID %d: %w", id, err)
	}
	return &user, nil
}

func FindByEmail(db *sql.DB, email string) (*User, error) {
	row := db.QueryRow("SELECT id, name, email FROM users WHERE email = ?", email)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("error scanning user with email %s: %w", email, err)
	}
	return &user, nil
}

func (u *User) Update(db *sql.DB) (int64, error) {
	stmt, err := db.Prepare("UPDATE users SET name = ?, email = ? WHERE id = ?")
	if err != nil {
		return 0, fmt.Errorf("error preparing update statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(u.Name, u.Email, u.ID)
	if err != nil {
		return 0, fmt.Errorf("error executing update statement: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error getting rows affected: %w", err)
	}
	log.Printf("User with ID %d updated. Rows affected: %d", u.ID, rowsAffected)
	return rowsAffected, nil
}

func (u *User) Delete(db *sql.DB) (int64, error) {
	stmt, err := db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		return 0, fmt.Errorf("error preparing delete statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(u.ID)
	if err != nil {
		return 0, fmt.Errorf("error executing delete statement: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error getting rows affected: %w", err)
	}
	log.Printf("User with ID %d deleted. Rows affected: %d", u.ID, rowsAffected)
	return rowsAffected, nil
}

func GetAllUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, fmt.Errorf("error querying all users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			log.Printf("Error scanning user row: %v", err)
			continue
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over user rows: %w", err)
	}
	return users, nil
}
