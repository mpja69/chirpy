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
	fdb.db.CreateTokenForUserId(user.Id, refreshToken, time.Hour*24*60)

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

// handleRevokeToken
// Takes no body, BUT the refresh token
//
//	Returns no body, BUT status 204
func (fdb *apiConfig) handleRevokeToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	err = fdb.db.RevokeToken(refreshToken)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// handleRefreshToken
// Takes no body, BUT the refresh token
//
//	Returns no body, BUT status 401 for invalid token, 200 for success
func (fdb *apiConfig) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	oldTokenString, err := auth.GetBearerToken(r)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	oldToken, err := fdb.db.GetToken(oldTokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	expirationTime := time.Unix(oldToken.ExpirationTime, 0)
	if expirationTime.Before(time.Now().UTC()) {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	newTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = fdb.db.RefreshToken(oldTokenString, newTokenString, time.Hour*24*60)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	type returnValue struct {
		Token string `json:"token"`
	}
	sendJsonResponse(w, http.StatusOK, returnValue{Token: newTokenString})

}
