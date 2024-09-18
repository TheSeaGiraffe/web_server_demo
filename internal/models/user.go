package models

import (
	"cmp"
	"errors"
	"slices"

	"golang.org/x/crypto/bcrypt"
)

var ErrUserNotExist = errors.New("User does not exist")

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

const CryptCost = 12

// Not sure if I even need this function. Will keep it for now.
func (db *DB) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), CryptCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Think of a better way of doing this later
func (db *DB) EmailExists(email string) (bool, error) {
	_, err := db.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, ErrUserNotExist) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Load db
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	// Check if user with specified email exists
	for _, user := range dbStruct.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, ErrUserNotExist
}

func (db *DB) GetUserByID(id int) (User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Load db
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	// Check if user with specified id exists
	user, ok := dbStruct.Users[id]
	if ok {
		return user, nil
	}

	return User{}, ErrUserNotExist
}

func (db *DB) CreateUser(email, password string) (User, error) {
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
	hashedPass, err := db.hashPassword(password)
	if err != nil {
		return User{}, err
	}

	lastID++
	user := User{
		ID:          lastID,
		Email:       email,
		Password:    hashedPass,
		IsChirpyRed: false,
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

func (db *DB) UpdateUser(id int, email, password string) error {
	// Lock db and defer unlocking
	db.mu.Lock()
	defer db.mu.Unlock()

	// Load db
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}

	// Check that user with matching ID actually exists
	user, ok := dbStruct.Users[id]
	if !ok {
		return ErrUserNotExist
	}

	// Write updated user info to disk
	hashedPass, err := db.hashPassword(password)
	if err != nil {
		return err
	}
	user = User{
		ID:          id,
		Email:       email,
		Password:    hashedPass,
		IsChirpyRed: user.IsChirpyRed,
	}
	dbStruct.Users[id] = user
	err = db.writeDB(dbStruct)
	if err != nil {
		return err
	}

	return nil
}
