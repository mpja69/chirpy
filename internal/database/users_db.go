package database

import (
	"fmt"
)

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUser(email string, password string) (User, error) {
	fmt.Printf("Creating user with email: %s\n", email)

	_, err := db.GetUserByEmail(email)
	if err == nil {
		return User{}, fmt.Errorf("User (%s) already exist", email)
	}

	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	id := len(dbs.Users) + 1
	fmt.Printf("Assigned ID: %d\n", id)
	user := User{
		Id:       id,
		Email:    email,
		Password: password,
	}

	dbs.SetUser(id, user)

	err = db.writeDB(dbs)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetUsers() ([]User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(dbs.Users))
	for _, user := range dbs.Users {
		users = append(users, user)
	}
	return users, nil
}

func (db *DB) GetUserById(id int) (User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbs.GetUser(id)

	if ok {
		return user, nil
	}
	return User{}, fmt.Errorf("User Id not in db")
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbs.Users {
		// fmt.Printf("Comparing users: %v - %s", user, email)
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, fmt.Errorf("User (email) not in db")
}
