package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) LikePhoto(photo_id string, user_id string) error {
	var owner_id string
	err := db.c.QueryRow("SELECT owner FROM Photo WHERE id == ?;", photo_id).Scan(&owner_id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrPhotoDoesNotExist
	} else if err != nil {
		return err
	} else if user_id == owner_id {
		return ErrNotSelfLike
	}

	var banned bool
	banned, err = db.CheckBan(owner_id, user_id)
	if err != nil {
		return err
	} else if banned {
		return ErrBanned
	}

	sqlStmt := `INSERT OR IGNORE INTO Like (photo, user) VALUES (?, ?)`
	_, err = db.c.Exec(sqlStmt, photo_id, user_id)
	return err
}
