package database

func (db *appdbimpl) UnlikePhoto(photo_id string, user_id string) error {
	exists, err := db.CheckPhoto(photo_id)
	if err != nil {
		return err
	} else if !exists {
		return ErrPhotoDoesNotExist
	}

	sqlStmt := `DELETE FROM Like WHERE photo == ? AND user == ?;`
	_, err = db.c.Exec(sqlStmt, photo_id, user_id)
	return err
}
