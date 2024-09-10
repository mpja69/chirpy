package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// NOTE:	handleLogin, POST /api/login, authenticates and logs in a user
//
//	Takes email, password (and optional expiration time in seconds)
//	Returns a JWT token (and user Id and email )
func (fdb *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := fdb.db.GetUserByEmail(params.Email)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	secondsInDay := 24 * 60 * 60
	expires := params.ExpiresInSeconds
	if expires == 0 || expires > secondsInDay {
		expires = secondsInDay
	}

	// NOTE: Andra kör inte med MewNumericDate(...), men med .Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Second * time.Duration(expires))),
		Subject:   strconv.Itoa(user.Id),
	})

	signedToken, err := token.SignedString(fdb.jwtSecret)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	type ResponseUser struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
		Token string `json:"token"`
	}
	responseVal := ResponseUser{
		Id:    user.Id,
		Email: user.Email,
		Token: signedToken,
	}
	sendJsonResponse(w, http.StatusOK, responseVal)
}
