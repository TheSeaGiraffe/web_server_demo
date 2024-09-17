package models

import (
	"cmp"
	"fmt"
	"os"
	"reflect"
	"slices"
	"testing"
)

// Re-write the tests to use golden files since the database is really just a JSON file

const testDB = "chirp_db-test.json"

var testChirps = []Chirp{
	{1, "The first chirp", 1},
	{2, "Another chirp", 2},
	{5, "That was some great mac 'n cheese we had last night", 3},
	{10, "Anyone else gotta deal with noisy neighbors. I'm losing sleep over here!", 4},
}

// Setup test DB and populate it with test cases
func dbSetup() func() {
	// Remove any existing test DB files
	_, err := os.ReadFile(testDB)
	if err == nil {
		err = os.Remove(testDB)
		if err != nil {
			panic(fmt.Errorf("error deleting existing DB file: %v", err))
		}
	}

	// Setup a new test DB
	chirpDB, err := NewDB(testDB)
	if err != nil {
		panic(fmt.Errorf("error creating DB file: %v", err))
	}

	// Populate the test DB with test cases
	chirps := make(map[int]Chirp)
	for _, chirp := range testChirps {
		chirps[chirp.ID] = chirp
	}

	chirpDBStruct := DBStructure{
		Chirps: chirps,
	}

	err = chirpDB.writeDB(chirpDBStruct)
	if err != nil {
		panic(fmt.Errorf("error writing test cases to DB: %v", err))
	}

	// Return teardown function which just deletes the test DB
	return func() {
		err = os.Remove(testDB)
		if err != nil {
			panic(fmt.Errorf("error deleting DB file: %v", err))
		}
	}
}

// Figure out why the teardown function isn't working
func TestMain(m *testing.M) {
	dbTeardown := dbSetup()
	defer dbTeardown()
	code := m.Run()
	os.Exit(code)
}

func TestChirpDB(t *testing.T) {
	chirpDB, err := NewDB(testDB)
	if err != nil {
		t.Fatalf("could not establish database connection: %v", err)
	}

	t.Run("GetChirps", testChirpDB_GetChirps(chirpDB))
	t.Run("CreateChirp", testChirpDB_CreateChirp(chirpDB))
}

func testChirpDB_GetChirps(chirpDB *DB) func(t *testing.T) {
	return func(t *testing.T) {
		chirps, err := chirpDB.GetChirps()
		if err != nil {
			t.Fatalf("could not retrieve chirps: %v", err)
		}
		slices.SortFunc(chirps, func(a, b Chirp) int {
			return cmp.Compare(a.ID, b.ID)
		})

		if !reflect.DeepEqual(chirps, testChirps) {
			t.Errorf("chirps != testChirps")
		}
	}
}

func testChirpDB_CreateChirp(chirpDB *DB) func(t *testing.T) {
	return func(t *testing.T) {
		newChirpBody := "Have you guys checked out that new pizza place yet?"
		chirp, err := chirpDB.CreateChirp(newChirpBody, 1)
		if err != nil {
			t.Fatalf("could not create chirp: %v", err)
		}

		chirpDBStruct, err := chirpDB.loadDB()
		if err != nil {
			t.Fatalf("could not load chirps from DB: %v", err)
		}

		dbChirp := chirpDBStruct.Chirps[11]

		if !reflect.DeepEqual(chirp, dbChirp) {
			t.Errorf("chirp != dbChirp")
		}
	}
}
