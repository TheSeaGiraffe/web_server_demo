package models

import (
	"cmp"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"slices"
	"time"
)

const (
	RefreshTokenLen   = 32
	TokenExpiryInDays = 60
)

var ErrTokenNotExist = errors.New("Token does not exist")

type Token struct {
	ID        int       `json:"id"`
	Plaintext string    `json:"plaintext"`
	Expiry    time.Time `json:"expiry"`
	UserID    int       `json:"user_id"`
}

func (db *DB) generateRefreshToken(n int) (string, error) {
	byteArr := make([]byte, n)

	_, err := rand.Read(byteArr)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(byteArr), nil
}

func (db *DB) tokenByUserID(refreshTokens []Token, userID int) (Token, error) {
	for _, token := range refreshTokens {
		if token.UserID == userID {
			return token, nil
		}
	}
	return Token{}, ErrTokenNotExist
}

func (db *DB) createNewToken(tokenPlaintext string, userID, lastID int) Token {
	return Token{
		ID:        lastID,
		Plaintext: tokenPlaintext,
		Expiry:    time.Now().Add(TokenExpiryInDays * 24 * time.Hour),
		UserID:    userID,
	}
}

func (db *DB) GetRefreshTokenByteLen(tokenPlaintext string) (int, error) {
	hexBytes, err := hex.DecodeString(tokenPlaintext)
	if err != nil {
		return 0, err
	}
	return len(hexBytes), nil
}

func (db *DB) CreateRefreshToken(userID int) (Token, error) {
	// Lock db and defer unlocking
	db.mu.Lock()
	defer db.mu.Unlock()

	// Load db
	dbStruct, err := db.loadDB()
	if err != nil {
		return Token{}, err
	}

	// Generate refresh token
	tokenPlaintext, err := db.generateRefreshToken(RefreshTokenLen)
	if err != nil {
		return Token{}, err
	}

	// Get tokens if they exist, sort them, and identify the latest ID
	lastID := 0
	var tokens []Token
	if len(dbStruct.Tokens) > 0 {
		for _, token := range dbStruct.Tokens {
			tokens = append(tokens, token)
		}
		slices.SortFunc(tokens, func(a, b Token) int {
			return -cmp.Compare(a.ID, b.ID)
		})
		lastID = tokens[0].ID
	}

	// Check if token exists for current user in the event that the user is logging in again
	// and then overwrite it.
	// Find a way to rewrite this in a way that eliminates the redundancy later.
	var token Token
	if len(tokens) > 0 {
		token, err = db.tokenByUserID(tokens, userID)
		if errors.Is(err, ErrTokenNotExist) {
			lastID++
			token = db.createNewToken(tokenPlaintext, userID, lastID)
		} else {
			lastID = token.ID
			token.Plaintext = tokenPlaintext
		}
	} else {
		lastID++
		token = db.createNewToken(tokenPlaintext, userID, lastID)
	}

	// Write token to disk
	if len(dbStruct.Tokens) == 0 {
		dbStruct.Tokens = make(map[int]Token)
	}
	dbStruct.Tokens[lastID] = token
	err = db.writeDB(dbStruct)
	if err != nil {
		return Token{}, err
	}

	return token, nil
}

// For revoking refresh tokens
func (db *DB) DeleteRefreshToken(tokenID int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}

	_, ok := dbStruct.Tokens[tokenID]
	if !ok {
		return fmt.Errorf("Token with ID '%d' does not exist", tokenID)
	}

	delete(dbStruct.Tokens, tokenID)

	err = db.writeDB(dbStruct)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) RefreshTokenExpired(tokenPlaintext string) (bool, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStruct, err := db.loadDB()
	if err != nil {
		return false, err
	}

	for _, token := range dbStruct.Tokens {
		if token.Plaintext == tokenPlaintext {
			log.Printf("Current time: %v", time.Now())
			log.Printf("Token expiry: %v", token.Expiry)
			log.Printf("Current time is after token expiry: %v", time.Now().After(token.Expiry))
			return time.Now().After(token.Expiry), nil
		}
	}

	return false, ErrTokenNotExist
}

func (db *DB) GetUserByRefreshToken(tokenPlaintext string) (User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	var userID int
	for _, token := range dbStruct.Tokens {
		if token.Plaintext == tokenPlaintext {
			userID = token.UserID
		}
	}

	if userID == 0 {
		return User{}, ErrTokenNotExist
	}

	user := dbStruct.Users[userID]

	return user, nil
}
