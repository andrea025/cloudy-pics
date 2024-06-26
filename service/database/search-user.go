package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) SearchUser(username string, req_id string) (UserShortInfo, error) {
	var user UserShortInfo
	err := db.c.QueryRow("SELECT * FROM User WHERE username == ?;", username).Scan(&user.Id, &user.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return user, ErrUserDoesNotExist
	} else if err != nil {
		return user, err
	}

	var banned bool
	banned, err = db.CheckBan(user.Id, req_id)
	if err != nil {
		return user, err
	} else if banned {
		return UserShortInfo{}, nil
	}

	return user, nil
}
