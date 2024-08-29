package database

import (
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"sync"
)

const DBFilePath = "chirp_db.json"

type DB struct {
	path string
	mu   sync.RWMutex
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
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

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	// Lock db and defer unlocking
	db.mu.Lock()
	defer db.mu.Unlock()

	// Load db
	chirpDBStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	// Get the last ID (i.e., the largest ID)
	var chirps []Chirp
	lastID := 0
	if len(chirpDBStruct.Chirps) > 0 {
		for _, chirp := range chirpDBStruct.Chirps {
			chirps = append(chirps, chirp)
		}

		// This should sort in descending order
		slices.SortFunc(chirps, func(a, b Chirp) int {
			return -cmp.Compare(a.ID, b.ID)
		})

		lastID = chirps[0].ID
	}

	// Create chirp
	lastID++
	chirp := Chirp{
		ID:   lastID,
		Body: body,
	}

	// Write chirp to disk
	if len(chirpDBStruct.Chirps) == 0 {
		chirpDBStruct.Chirps = make(map[int]Chirp)
	}
	chirpDBStruct.Chirps[lastID] = chirp
	err = db.writeDB(chirpDBStruct)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	chirpDBStruct, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}

	if len(chirpDBStruct.Chirps) == 0 {
		return []Chirp{}, nil
	}

	var chirps []Chirp
	for _, chirp := range chirpDBStruct.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirpByID(id int) (Chirp, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	chirpDBStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	if len(chirpDBStruct.Chirps) == 0 {
		return Chirp{}, fmt.Errorf("No chirps in database")
	}

	chirp, ok := chirpDBStruct.Chirps[id]
	if !ok {
		return Chirp{}, fmt.Errorf("Chirp with ID '%d' does not exist", id)
	}

	return chirp, nil
}
