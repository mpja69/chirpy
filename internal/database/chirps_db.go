package database

import "fmt"

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	// db.mux.Lock()
	// defer db.mux.Unlock()
	dbs, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	id := len(dbs.Chirps) + 1
	chirp := Chirp{
		Id:   id,
		Body: body,
	}

	dbs.SetChirp(id, chirp)

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
		return nil, err
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

	chirp, ok := dbs.GetChirp(id)

	if ok {
		return chirp, nil
	}
	return Chirp{}, fmt.Errorf("Chirp Id not in db")
}
