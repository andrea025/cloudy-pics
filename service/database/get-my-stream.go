package database

func (db *appdbimpl) GetMyStream(id string) ([]Photo, error) {
	stream := []Photo{}

	sqlStmt := `SELECT id FROM Photo WHERE owner in (select user_followed from Following where user_following = ? AND user_followed not in (select user_banning from Banned where user_banned = ?)) ORDER BY created_at DESC LIMIT 50;`
	rows, err := db.c.Query(sqlStmt, id, id)
	if err != nil {
		return []Photo{}, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var pid string
		err = rows.Scan(&pid)
		if err != nil {
			return []Photo{}, err
		}

		var photo Photo
		photo, err = db.GetPhoto(pid, id)
		if err != nil {
			return []Photo{}, err
		}
		stream = append(stream, photo)
	}
	if err = rows.Err(); err != nil {
		return []Photo{}, err
	}

	return stream, nil
}
