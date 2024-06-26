package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) SetMyUsername(id string, username string) error {
	var uid string
	err := db.c.QueryRow("SELECT id FROM User WHERE id <> ? AND username == ?;", id, username).Scan(&uid)
	if errors.Is(err, sql.ErrNoRows) {
		// no other user has chosen the new username
		sqlStmt := `UPDATE User SET username = ? WHERE id == ?;`
		_, err = db.c.Exec(sqlStmt, username, id)
		return err
	} else if err != nil {
		return err
	}

	return ErrUsernameAlreadyTaken
}
