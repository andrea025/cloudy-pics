package database

import (
	"database/sql"
	"errors"
)

type Comment struct {
	Id              string
	Photo           string
	User            UserShortInfo
	Text            string
	CreatedDatetime string
}

func (db *appdbimpl) CommentPhoto(cid string, pid string, uid string, text string, created_datetime string) (Comment, error) {
	var owner_id string
	err := db.c.QueryRow("SELECT owner FROM Photo WHERE id == ?;", pid).Scan(&owner_id)
	if errors.Is(err, sql.ErrNoRows) {
		return Comment{}, ErrPhotoDoesNotExist
	} else if err != nil {
		return Comment{}, err
	}

	var user UserShortInfo
	err = db.c.QueryRow("SELECT id, username FROM User WHERE id == ?;", uid).Scan(&user.Id, &user.Username)
	if err != nil {
		return Comment{}, err
	}

	var banned bool
	banned, err = db.CheckBan(owner_id, uid)
	if err != nil {
		return Comment{}, err
	} else if banned {
		return Comment{}, ErrBanned
	}

	comment := Comment{Id: cid, Photo: pid, User: user, CreatedDatetime: created_datetime, Text: text}

	sqlStmt := `INSERT INTO Comment (id, photo, user, text, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err = db.c.Exec(sqlStmt, comment.Id, comment.Photo, comment.User.Id, comment.Text, comment.CreatedDatetime)
	return comment, err
}
