package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetPhoto(id string, req_id string) (Photo, error) {
	var photo Photo
	var owner_id string

	err := db.c.QueryRow("SELECT * FROM Photo WHERE id=?", id).Scan(&photo.Id, &photo.CreatedDatetime, &photo.PhotoUrl, &owner_id)
	if errors.Is(err, sql.ErrNoRows) {
		return photo, ErrPhotoDoesNotExist
	}

	// check if the requesting user has been banned
	var banned bool
	banned, err = db.CheckBan(owner_id, req_id)
	if err != nil {
		return photo, err
	} else if banned {
		return photo, ErrBanned
	}

	err = db.c.QueryRow("SELECT * FROM User WHERE id=?", owner_id).Scan(&photo.Owner.Id, &photo.Owner.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return photo, ErrPhotoDoesNotExist
	}

	photo.Likes = LikesCollection{Count: 0, Users: []UserShortInfo{}}
	rows, er := db.c.Query("SELECT id, username FROM User, Like WHERE id=user AND photo=? LIMIT 50;", id)
	if er != nil {
		return photo, er
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var user UserShortInfo
		er = rows.Scan(&user.Id, &user.Username)
		if er != nil {
			return photo, er
		}

		var banned bool
		banned, er = db.CheckBan(user.Id, req_id)
		if er != nil {
			return photo, err
		} else if !banned {
			photo.Likes.Users = append(photo.Likes.Users, user)
		}
		photo.Likes.Count++
	}
	if er = rows.Err(); er != nil {
		return photo, er
	}

	photo.Comments = CommentsCollection{Count: 0, Comments: []Comment{}}
	rows, er = db.c.Query("SELECT * FROM Comment WHERE photo=? LIMIT 50;", id)
	if er != nil {
		return photo, er
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var comment Comment
		var uid string
		er = rows.Scan(&comment.Id, &comment.Photo, &uid, &comment.Text, &comment.CreatedDatetime)
		if er != nil {
			return photo, er
		}

		er = db.c.QueryRow("SELECT id, username FROM User WHERE id == ?;", uid).Scan(&comment.User.Id, &comment.User.Username)
		if er != nil {
			return photo, er
		}

		var banned bool
		banned, er = db.CheckBan(comment.User.Id, req_id)
		if er != nil {
			return photo, er
		} else if !banned {
			photo.Comments.Comments = append(photo.Comments.Comments, comment)
		}
		photo.Comments.Count++
	}
	if er = rows.Err(); er != nil {
		return photo, er
	}

	return photo, nil
}
