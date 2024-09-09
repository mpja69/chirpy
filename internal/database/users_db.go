package database

import (
	"fmt"
)

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUser(email string, password string) (User, error) {
	fmt.Printf("Creating user with email: %s\n", email)

	_, err := db.GetUserByEmail(email)
	if err == nil {
		return User{}, ErrDuplicate
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

// CreateUser creates a new user and saves it to disk
func (db *DB) UpdateUser(id int, email string, password string) (User, error) {
	fmt.Printf("Updating user: %d\n", id)

	// INFO: Check that the user exixsts, so we can update it
	user, err := db.GetUserById(id)
	if err != nil {
		return User{}, err
	}

	// INFO: Check if the email changed...and if so, is the new email already in db
	if email != user.Email {
		_, err = db.GetUserByEmail(email)
		if err == nil {
			return User{}, ErrDuplicate
		}
	}

	// INFO: Set the new email and pw...even if some might be the same as before
	user.Email = email
	user.Password = password

	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
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
	return User{}, ErrNotExist
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbs.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, ErrNotExist
}
