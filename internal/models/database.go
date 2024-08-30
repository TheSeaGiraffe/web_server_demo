package models

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

const DBFilePath = "chirp_db.json"

type DB struct {
	path string
	mu   sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

// NewDB creates a new database connection and creates a database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	chirpDB := &DB{
		path: path,
	}
	err := chirpDB.ensureDB()
	if err != nil {
		return nil, err
	}

	return chirpDB, nil
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if os.IsNotExist(err) {
		err = os.WriteFile(db.path, []byte{}, 0644)
		if err != nil {
			return fmt.Errorf("could not create DB file: %w", err)
		}
	}

	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	dbFile, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, fmt.Errorf("error reading DB file: %w", err)
	}

	if len(dbFile) == 0 {
		// This should only happen if the database is empty which is only the case
		// if you run the server without a DB. An empty DB is still a valid state
		// and shouldn't error
		return DBStructure{}, nil
	}

	var chirpDBStruct DBStructure
	err = json.Unmarshal(dbFile, &chirpDBStruct)
	if err != nil {
		return DBStructure{}, fmt.Errorf("error loading DB file: %w", err)
	}
	return chirpDBStruct, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	chirpsData, err := json.Marshal(dbStructure)
	if err != nil {
		return fmt.Errorf("error marshaling data: %w", err)
	}

	err = os.WriteFile(db.path, chirpsData, 0644)
	if err != nil {
		return fmt.Errorf("error writing to DB: %w", err)
	}

	return nil
}
