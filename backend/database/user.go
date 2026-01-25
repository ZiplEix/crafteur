package database

import (
	"database/sql"

	"github.com/ZiplEix/crafteur/core"
)

func CreateUser(u *core.User) error {
	_, err := DB.Exec("INSERT INTO users (id, username, password_hash) VALUES (?, ?, ?)", u.ID, u.Username, u.PasswordHash)
	return err
}

func GetUserByUsername(username string) (*core.User, error) {
	row := DB.QueryRow("SELECT id, username, password_hash FROM users WHERE username = ?", username)
	var u core.User
	if err := row.Scan(&u.ID, &u.Username, &u.PasswordHash); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}
	return &u, nil
}

func GetUserByID(id string) (*core.User, error) {
	row := DB.QueryRow("SELECT id, username, password_hash FROM users WHERE id = ?", id)
	var u core.User
	if err := row.Scan(&u.ID, &u.Username, &u.PasswordHash); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}
	return &u, nil
}
