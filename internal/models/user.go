package models

import (
	"cmp"
	"slices"
)

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (db *DB) CreateUser(email string) (User, error) {
	// Lock db and defer unlocking
	db.mu.Lock()
	defer db.mu.Unlock()

	// Load db
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	// Get the last ID (i.e., the largest ID)
	var users []User
	lastID := 0
	if len(dbStruct.Users) > 0 {
		for _, user := range dbStruct.Users {
			users = append(users, user)
		}

		// This should sort in descending order
		slices.SortFunc(users, func(a, b User) int {
			return -cmp.Compare(a.ID, b.ID)
		})

		lastID = users[0].ID
	}

	// Create user
	lastID++
	user := User{
		ID:    lastID,
		Email: email,
	}

	// Write user to disk
	if len(dbStruct.Users) == 0 {
		dbStruct.Users = make(map[int]User)
	}
	dbStruct.Users[lastID] = user
	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
