package database

func (db *appdbimpl) UnbanUser(user_id string, target_user_id string) error {
	exists, err := db.CheckUser(target_user_id)
	if err != nil {
		return err
	} else if !exists {
		return ErrUserDoesNotExist
	}

	sqlStmt := `DELETE FROM Banned WHERE user_banning == ? AND user_banned == ?`
	_, err = db.c.Exec(sqlStmt, user_id, target_user_id)
	return err
}
