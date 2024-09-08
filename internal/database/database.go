package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	err := db.ensureDB()
	return db, err
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	id := len(dbs.Chirps) + 1
	chirp := Chirp{
		Id:   id,
		Body: body,
	}
	dbs.Chirps[id] = chirp

	err = db.writeDB(dbs)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return nil, fmt.Errorf("GetChirps: %v", err)
	}

	chirps := make([]Chirp, 0, len(dbs.Chirps))
	for _, chirp := range dbs.Chirps {
		chirps = append(chirps, chirp)
	}
	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	chirp, ok := dbs.Chirps[id]
	if ok {
		return chirp, nil
	}
	return Chirp{}, fmt.Errorf("Id not in db")
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.Stat(db.path)
	if err == nil {
		return nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}

	return err
}

func (db *DB) createDB() error {
	dbs := DBStructure{
		Chirps: map[int]Chirp{},
	}
	return db.writeDB(dbs)
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbs := DBStructure{}
	data, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbs, nil
	}

	err = json.Unmarshal(data, &dbs)
	if err != nil {
		return dbs, fmt.Errorf("LoadDB: %v", err)
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

	err = os.WriteFile(db.path, data, 0600)
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

// func (db *DB) WrapperHandlerFunc(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		next(w, r)
// 	})
// }

// func (db *DB) WrapperHandler(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		next.ServeHTTP(w, r)
// 	})
// }
