package database

import (
	"fmt"
	"time"
)

type RefreshToken struct {
	UserId    int       `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	Token     string    `json:"token"`
}

// SaveTokenForUserId - Saves a refresh token, connected to a UserId and expiration time
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
		UserId:    userId,
		ExpiresAt: time.Now().Add(time.Hour * duration),
		Token:     tokenString,
	}

	dbs.SetToken(tokenString, token)

	err = db.writeDB(dbs)
	if err != nil {
		return RefreshToken{}, err
	}

	return token, nil
}

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

// GetToke - nLoad file, Get value, Return value
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
