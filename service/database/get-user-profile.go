package database

import (
	"database/sql"
	"errors"
)

type User struct {
	Id             string
	Username       string
	Followers      int
	Following      int
	UploadedPhotos int
	Photos         []Photo
}

func (db *appdbimpl) GetUserProfile(id string, req_id string) (User, error) {
	var user User

	// check if the user exists
	err := db.c.QueryRow("SELECT * FROM User WHERE id=?", id).Scan(&user.Id, &user.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return user, ErrUserDoesNotExist
	}

	// check if the requesting user has been banned
	var banned bool
	banned, err = db.CheckBan(id, req_id)
	if err != nil {
		return user, err
	} else if banned {
		return user, ErrBanned
	}

	// get the other values
	// ignore the error because the user id has already been validated and COUNT() always returns a value, even if 0
	_ = db.c.QueryRow("SELECT COUNT(*) FROM Following WHERE user_followed=?", id).Scan(&user.Followers)
	_ = db.c.QueryRow("SELECT COUNT(*) FROM Following WHERE user_following=?", id).Scan(&user.Following)
	_ = db.c.QueryRow("SELECT COUNT(*) FROM Photo WHERE owner=?", id).Scan(&user.UploadedPhotos)

	user.Photos = []Photo{}

	sqlStmt := `SELECT id FROM Photo WHERE owner == ?;`
	rows, er := db.c.Query(sqlStmt, id)
	if er != nil {
		return user, er
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var pid string
		er = rows.Scan(&pid)
		if er != nil {
			return user, er
		}

		var photo Photo
		photo, er = db.GetPhoto(pid, id)
		if er != nil {
			return user, er
		}
		user.Photos = append(user.Photos, photo)
	}
	if er = rows.Err(); er != nil {
		return user, er
	}

	return user, nil
}
