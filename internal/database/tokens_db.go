package database

import (
	"fmt"
	"time"
)

type RefreshToken struct {
	UserId         int       `json:"user_id"`
	ExpirationTime time.Time `json:"expiration_time"`
	Token          string    `json:"token"`
}

// SaveTokenForUserId Saves a refresh token
func (db *DB) SaveTokenForUserId(userId int, tokenString string, duration time.Duration) (RefreshToken, error) {
	fmt.Printf("Creating token: %s, for userId: %d, at time: %v \n", tokenString, userId, duration)

	//TODO: Kolla om det redan finns ett token för userId
	// _, err := db.GetTokenByUserId(userId)
	// if err == nil {
	// 	return Token{}, ErrDuplicate
	// }

	dbs, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}
	token := RefreshToken{
		UserId:         userId,
		ExpirationTime: time.Now().Add(time.Hour * duration),
		Token:          tokenString,
	}

	dbs.SetToken(tokenString, token)

	err = db.writeDB(dbs)
	if err != nil {
		return RefreshToken{}, err
	}

	return token, nil
}

// WARN: Not needed
// GetTokenByUserId - Helper function to avoid duplicates
// func (db *DB) GetTokenByUserId(userId int) (Token, error) {
// 	fmt.Printf("Searching for token belonging to userId: %d\n", userId)
// 	dbs, err := db.loadDB()
// 	if err != nil {
// 		return Token{}, err
// 	}
//
// 	for _, token := range dbs.Tokens {
// 		if token.UserId == userId {
// 			return token, nil
// 		}
// 	}
// 	return Token{}, ErrNotExist
// }

// Revoke token - Just delete from map and file
func (db *DB) RevokeToken(tokenString string) error {
	fmt.Printf("Revoking (deleting) token: %s\n", tokenString)

	dbs, err := db.loadDB()
	if err != nil {
		return err
	}

	dbs.DeleteToken(tokenString)

	err = db.writeDB(dbs)
	if err != nil {
		return err
	}

	return nil
}

// RefreshToken - Find old token, Copy userID, Set new expiration time, Save the new and delete the old
func (db *DB) RefreshToken(oldTokenString, newTokenString string, duration time.Duration) error {
	dbs, err := db.loadDB()
	if err != nil {
		return err
	}
	oldtoken, ok := dbs.GetToken(oldTokenString)
	if !ok {
		return ErrNotExist
	}
	// dbs.DeleteToken(oldTokenString)

	newToken := RefreshToken{
		UserId:         oldtoken.UserId,
		ExpirationTime: time.Now().Add(time.Hour * duration),
		Token:          newTokenString,
	}

	dbs.SetToken(newTokenString, newToken)

	err = db.writeDB(dbs)
	if err != nil {
		return err
	}
	return nil
}

// Load file, Get value, Return value
func (db *DB) GetToken(tokenString string) (RefreshToken, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}

	token, ok := dbs.GetToken(tokenString)
	if ok {
		return token, nil
	}
	return RefreshToken{}, ErrNotExist
}

// WARN: Not used
// Load file, Set value, Save file
// func (db *DB) SetToken(tokenString string, token Token) error {
// 	dbs, err := db.loadDB()
// 	if err != nil {
// 		return err
// 	}
// 	dbs.SetToken(tokenString, token)
// 	return nil
// }

// WARN: Not used
// Load file, Delete value, Save file
// func (db *DB) DeleteToken(tokenString string) error {
// 	dbs, err := db.loadDB()
// 	if err != nil {
// 		return err
// 	}
// 	dbs.DeleteToken(tokenString)
// 	return nil
// }
