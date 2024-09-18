package models

type EventChirpyRed struct {
	Event string `json:"event"`
	Data  struct {
		UserID int `json:"user_id"`
	} `json:"data"`
}

func (db *DB) UpgradeChirpyRedForUser(userID int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}

	user, ok := dbStruct.Users[userID]
	if !ok {
		return ErrUserNotExist
	}

	// Wonder if there is a better way of doing this
	user.IsChirpyRed = true
	dbStruct.Users[userID] = user

	err = db.writeDB(dbStruct)
	if err != nil {
		return err
	}

	return nil
}
