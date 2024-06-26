package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) DeletePhoto(photo_id string, req_id string) (string, error) {
	var url, owner string
	err := db.c.QueryRow("SELECT url, owner FROM Photo WHERE id=?", photo_id).Scan(&url, &owner)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrPhotoDoesNotExist
	} else if err != nil {
		return "", err
	}

	if owner != req_id {
		return "", ErrDeletePhotoForbidden
	}

	sqlStmt := `DELETE FROM Photo WHERE id == ?;`
	_, err = db.c.Exec(sqlStmt, photo_id)
	return url, err
}
