package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) CheckPhoto(photo_id string) (bool, error) {
	var id string
	err := db.c.QueryRow("SELECT id FROM Photo WHERE id == ?;", photo_id).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		// considering just the error, false will be ignored
		return false, err
	}
	return true, nil
}
