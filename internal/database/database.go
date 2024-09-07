package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

func (db *DB) WrapperHandlerFunc(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Do something...
		next(w, r)
		// TODO: Do something...
	})
}

// func (db *DB) WrapperHandler(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		next.ServeHTTP(w, r)
// 	})
// }

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	err := db.ensureDB()
	if err != nil {
		return nil, err
	}

	return &db, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp := Chirp{
		Id:   len(dbs.Chirps),
		Body: body,
	}
	dbs.Chirps[len(dbs.Chirps)] = chirp

	db.writeDB(dbs)
	return chirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {

	dbs, err := db.loadDB()
	if err != nil {
		return nil, fmt.Errorf("GetChirps: %v", err)
	}

	var chirps []Chirp
	for _, chirp := range dbs.Chirps {
		chirps = append(chirps, chirp)
	}
	return chirps, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.Stat(db.path)
	if err == nil {
		return nil
	}

	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(db.path)
		defer f.Close()
		if err != nil {
			return err
		}
		return nil
	}

	return err
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	data, err := os.ReadFile(db.path)

	if err != nil {
		return DBStructure{}, fmt.Errorf("LoadDB: %v", err)
	}
	if len(data) == 0 {
		return DBStructure{Chirps: make(map[int]Chirp)}, nil
	}

	var dbs DBStructure
	err = json.Unmarshal(data, &dbs)
	if err != nil {
		return DBStructure{}, fmt.Errorf("LoadDB: %v", err)
	}
	return dbs, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()
	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, data, fs.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}
type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}
