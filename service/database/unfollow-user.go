package database

func (db *appdbimpl) UnfollowUser(user_id string, target_user_id string) error {
	exists, err := db.CheckUser(target_user_id)
	if err != nil {
		return err
	} else if !exists {
		return ErrUserDoesNotExist
	}

	sqlStmt := `DELETE FROM Following WHERE user_following == ? AND user_followed == ?`
	_, err = db.c.Exec(sqlStmt, user_id, target_user_id)
	return err
}
