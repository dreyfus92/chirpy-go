package database

import (
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
