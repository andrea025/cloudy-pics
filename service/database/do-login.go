package database

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) DoLogin(username string) (string, error) {
	var id string
	err := db.c.QueryRow("SELECT id FROM User WHERE username=?", username).Scan(&id)

	if errors.Is(err, sql.ErrNoRows) {
		// if there is no user in the User table with the given username
		// generate a new user id
		id = fmt.Sprintf("%x", md5.Sum([]byte(username)))

		sqlStmt := `INSERT INTO User (id, username) VALUES (?, ?)`
		_, err = db.c.Exec(sqlStmt, id, username)
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}

	return id, nil
}
