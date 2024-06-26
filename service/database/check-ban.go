package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) CheckBan(user_id string, target_user_id string) (bool, error) {
	var uid, tuid string
	err := db.c.QueryRow("SELECT * FROM Banned WHERE user_banning == ? AND user_banned == ?;", user_id, target_user_id).Scan(&uid, &tuid)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		// considering just the error, false will be ignored
		return false, err
	}
	return true, nil
}
