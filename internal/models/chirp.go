package models

import (
	"cmp"
	"fmt"
	"slices"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	// Lock db and defer unlocking
	db.mu.Lock()
	defer db.mu.Unlock()

	// Load db
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	// Get the last ID (i.e., the largest ID)
	var chirps []Chirp
	lastID := 0
	if len(dbStruct.Chirps) > 0 {
		for _, chirp := range dbStruct.Chirps {
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
	if len(dbStruct.Chirps) == 0 {
		dbStruct.Chirps = make(map[int]Chirp)
	}
	dbStruct.Chirps[lastID] = chirp
	err = db.writeDB(dbStruct)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStruct, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}

	if len(dbStruct.Chirps) == 0 {
		return []Chirp{}, nil
	}

	var chirps []Chirp
	for _, chirp := range dbStruct.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirpByID(id int) (Chirp, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	if len(dbStruct.Chirps) == 0 {
		return Chirp{}, fmt.Errorf("No chirps in database")
	}

	chirp, ok := dbStruct.Chirps[id]
	if !ok {
		return Chirp{}, fmt.Errorf("Chirp with ID '%d' does not exist", id)
	}

	return chirp, nil
}
