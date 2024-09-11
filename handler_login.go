package main

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/mpja69/chirpy/internal/auth"
)

// NOTE:	handleLogin, POST /api/login, authenticates and logs in a user
//
//	Takes email, password
//	Returns userId, email, access token (JWT), refresh token
func (fdb *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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

	accessToken, err := auth.MakeJWT(user.Id, fdb.jwtSecret, time.Hour)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// TODO: Spara refresh token och dess exp-time i databasen. Egen tabell?
	expirationTime := time.Now().Add(time.Hour)
	fdb.db.CreateTokenForUserId(user.Id, refreshToken, expirationTime)

	type ResponseUser struct {
		Id           int    `json:"id"`
		Email        string `json:"email"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	responseVal := ResponseUser{
		Id:           user.Id,
		Email:        user.Email,
		Token:        accessToken,
		RefreshToken: refreshToken,
	}
	sendJsonResponse(w, http.StatusOK, responseVal)
}
