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
	path     string
	mux      *sync.RWMutex
	dbstruct *DBStructure
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

/*
	database functions:
	- NewDB
	- loadDB
	- writeDB
	- GetChirp
	- GetChirps
	- CreateChirp
	- CreateUser
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

	dbStruct := DBStructure{
		Chirps: make(map[int]Chirp),
		Users:  make(map[int]User),
	}

	// create a new DB struct
	db := DB{
		path:     path,
		mux:      &sync.RWMutex{},
		dbstruct: &dbStruct,
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

	// Decode JSON into db.dbstruct struct
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&db.dbstruct); err != nil {
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
	if err := encoder.Encode(db.dbstruct.Chirps); err != nil {
		return err
	}

	return nil
}

// GetChirp returns a SINGLE chirp from the database when given an ID
func (db *DB) GetChirp(id int) (Chirp, error) {
	// lock for readers
	db.mux.RLock()
	defer db.mux.RUnlock()

	chirp, ok := db.dbstruct.Chirps[id]

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
	for key := range db.dbstruct.Chirps {
		chirps = append(chirps, db.dbstruct.Chirps[key])
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
	for key := range db.dbstruct.Chirps {
		if key > maxKey {
			maxKey = key
		}
	}
	newKey := maxKey + 1

	// create the chirp and add it to the db
	db.dbstruct.Chirps[newKey] = Chirp{
		Id:   newKey,
		Body: body,
	}

	// write the db to disk
	db.writeDB()

	// return the new chirp
	return newKey
}

// Creating a new user
func (db *DB) CreateUser(email string) int {
	// lock for writers
	db.mux.Lock()
	defer db.mux.Unlock()

	// get new ID
	maxKey := 0
	for key := range db.dbstruct.Users {
		if key > maxKey {
			maxKey = key
		}
	}
	newKey := maxKey + 1

	// create the user and add it to the db
	db.dbstruct.Users[newKey] = User{
		Id:    newKey,
		Email: email,
	}

	// write the db to disk
	db.writeDB()

	// return the new user
	return newKey
}
