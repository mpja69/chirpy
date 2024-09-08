package database

import (
	"encoding/json"
	"errors"
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
		Users:  map[int]User{},
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
		return dbs, err
	}

	err = json.Unmarshal(data, &dbs)
	if err != nil {
		return dbs, err
	}

	dbs.mux = &sync.Mutex{}
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
	Users  map[int]User  `json:"users"`
	mux    *sync.Mutex
}
type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}
type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func (dbs *DBStructure) GetUser(id int) (User, bool) {
	dbs.mux.Lock()
	defer dbs.mux.Unlock()
	user, ok := dbs.Users[id]
	return user, ok
}

func (dbs *DBStructure) GetChirp(id int) (Chirp, bool) {
	dbs.mux.Lock()
	defer dbs.mux.Unlock()
	chirp, ok := dbs.Chirps[id]
	return chirp, ok
}

func (dbs *DBStructure) SetUser(id int, user User) {
	dbs.mux.Lock()
	defer dbs.mux.Unlock()
	dbs.Users[id] = user
}

func (dbs *DBStructure) SetChirp(id int, chirp Chirp) {
	dbs.mux.Lock()
	defer dbs.mux.Unlock()
	dbs.Chirps[id] = chirp
}
