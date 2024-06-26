package database

func (db *appdbimpl) GetFollowing(id string, req_id string) ([]UserShortInfo, error) {
	users := []UserShortInfo{}

	exists, err := db.CheckUser(id)
	if err != nil {
		return []UserShortInfo{}, err
	} else if !exists {
		return []UserShortInfo{}, ErrUserDoesNotExist
	}

	// check if the requesting user has been banned
	var banned bool
	banned, err = db.CheckBan(id, req_id)
	if err != nil {
		return []UserShortInfo{}, err
	} else if banned {
		return []UserShortInfo{}, ErrBanned
	}

	sqlStmt := `SELECT user_followed, username FROM User, Following WHERE user_following == ? AND user_followed == id LIMIT 50;`
	rows, er := db.c.Query(sqlStmt, id)
	if er != nil {
		return []UserShortInfo{}, er
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var user UserShortInfo
		err = rows.Scan(&user.Id, &user.Username)
		if err != nil {
			return []UserShortInfo{}, err
		}

		var banned bool
		banned, er = db.CheckBan(user.Id, req_id)
		if er != nil {
			return []UserShortInfo{}, err
		} else if !banned {
			users = append(users, user)
		}
	}
	if err = rows.Err(); err != nil {
		return []UserShortInfo{}, err
	}

	return users, nil
}
