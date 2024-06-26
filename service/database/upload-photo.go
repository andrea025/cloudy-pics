package database

type Photo struct {
	Id              string
	CreatedDatetime string
	PhotoUrl        string
	Owner           UserShortInfo
	Likes           LikesCollection
	Comments        CommentsCollection
}

type UserShortInfo struct {
	Id       string
	Username string
}

type LikesCollection struct {
	Count int
	Users []UserShortInfo
}

type CommentsCollection struct {
	Count    int
	Comments []Comment
}

func (db *appdbimpl) UploadPhoto(id string, created_at string, url string, owner string) (Photo, error) {
	var photo Photo
	var user UserShortInfo
	err := db.c.QueryRow("SELECT * FROM User WHERE id=?", owner).Scan(&user.Id, &user.Username)
	if err != nil {
		return photo, err
	}
	photo.Id, photo.CreatedDatetime, photo.PhotoUrl, photo.Owner, photo.Likes, photo.Comments = id, created_at, url, user, LikesCollection{Count: 0, Users: []UserShortInfo{}}, CommentsCollection{Count: 0, Comments: []Comment{}}

	sqlStmt := `INSERT INTO Photo (id, created_at, url, owner) VALUES (?, ?, ?, ?)`
	_, err = db.c.Exec(sqlStmt, id, created_at, url, owner)
	return photo, err
}
