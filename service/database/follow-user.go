package database

func (db *appdbimpl) FollowUser(user_id string, target_user_id string) error {
	exists, err := db.CheckUser(target_user_id)
	if err != nil {
		return err
	} else if !exists {
		return ErrUserDoesNotExist
	}

	sqlStmt := `INSERT OR IGNORE INTO Following (user_following, user_followed) VALUES (?, ?)`
	_, err = db.c.Exec(sqlStmt, user_id, target_user_id)
	return err
}
