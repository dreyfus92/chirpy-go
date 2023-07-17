package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
)

type DB struct {
	path   string
	mux    *sync.RWMutex
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

/*
	DB Functions:
	- NewDB
	- loadDB
	- writeDB
	- GetChirp
	- GetChirps
	- CreateChirp
*/

// NewDB creates a new database connection
// and creates the database file if it doesn't exist

func NewDB(path string) (*DB, error) {
	// If the json doesn't exist, create it, or append to the file
	_, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// create new DB struct
	db := DB{
		path:   path,
		mux:    &sync.RWMutex{},
		Chirps: make(map[int]Chirp),
	}

	// load the JSON file contents into memory
	db.loadDB()

	return &db, nil
}

// loadDB reads the database file into memory

func (db *DB) loadDB() error {
	// lock for readers
	db.mux.RLock()
	defer db.mux.RUnlock()

	//get the JSON from the file, decode into the db.Chirps struct
	// Open file for reading
	file, err := os.Open(db.path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode JSON into DB.Chirps struct
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&db.Chirps); err != nil {
		return err
	}

	return nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB() error {
	// get the JSON file contents
	file, err := os.OpenFile(db.path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// write the Chirps map to the JSON file

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(db.Chirps); err != nil {
		return err
	}

	return nil
}

// GetChirp returns a SINGLE chirp from the database when given an ID
func (db *DB) GetChirp(id int) (Chirp, error) {
	// lock for readers
	db.mux.RLock()
	defer db.mux.RUnlock()

	chirp, ok := db.Chirps[id]

	if !ok {
		return Chirp{}, fmt.Errorf("chirp with id %d not found", id)
	}

	return chirp, nil
}

//GetChirps returns all chirps from the database

func (db *DB) GetChirps() ([]Chirp, error) {
	// lock for readers
	db.mux.RLock()
	defer db.mux.RUnlock()

	//get all chirps
	chirps := []Chirp{}
	for key := range db.Chirps {
		chirps = append(chirps, db.Chirps[key])
	}

	// sort slice of Chirp objects by ID
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})

	return chirps, nil
}

// CreateChirp creates a new chirp and saves it to disk

func (db *DB) CreateChirp(body string) int {
	// lock for writers
	db.mux.Lock()
	defer db.mux.Unlock()

	// get next ID
	maxKey := 0
	for key := range db.Chirps {
		if key > maxKey {
			maxKey = key
		}
	}
	newKey := maxKey + 1

	// create the chirp and add it to the db
	db.Chirps[newKey] = Chirp{
		Id:   newKey,
		Body: body,
	}

	// write the db to disk
	db.writeDB()

	// return the new chirp
	return newKey
}
