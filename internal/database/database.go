package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	Id   string `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	err := db.ensureDB()

	if err != nil {
		return nil, err
	}

	return db, nil
}

// ensureDB creates a new database file if it doesn't exist

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)

	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Creating new database at %s\n", db.path)
			err = os.WriteFile(db.path, []byte(`{"chirps":{}}`), 0644)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// loadDB reads the database file into memory

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dat, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	dbStruct := DBStructure{}
	err = json.Unmarshal(dat, &dbStruct)

	if err != nil {
		return DBStructure{}, err
	}

	return dbStruct, nil

}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	json, err := json.Marshal(dbStructure)

	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, json, 0644)
	if err != nil {
		return err
	}

	return nil
}

//GetChirps returns all chirps from the database

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStruct, err := db.loadDB()

	if err != nil {
		return nil, err
	}

	chirps := []Chirp{}

	for _, chirp := range dbStruct.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

// CreateChirp creates a new chirp and saves it to disk

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStruct, err := db.loadDB()

	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStruct.Chirps) + 1
	chirp := Chirp{
		Id:   fmt.Sprintf("%d", id),
		Body: body,
	}

	dbStruct.Chirps[id] = chirp

	err = db.writeDB(dbStruct)

	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}
