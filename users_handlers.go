package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// GET /api/users
func (fdb *apiConfig) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := fdb.db.GetUsers()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	//TODO: 4) Skicka tillbaka allt utom passwword
	type ResponseUser struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}
	sortedUsers := []ResponseUser{}
	for _, u := range users {
		sortedUsers = append(sortedUsers, ResponseUser{
			Id:    u.Id,
			Email: u.Email,
		})
	}
	sort.Slice(sortedUsers, func(i int, j int) bool {
		return sortedUsers[i].Id < sortedUsers[j].Id
	})

	sendJsonResponse(w, http.StatusOK, sortedUsers)
}

// GET /api/users/{userId}
func (cfg *apiConfig) handleGetUserById(w http.ResponseWriter, r *http.Request) {
	idValue := r.PathValue("userId")
	id, err := strconv.Atoi(idValue)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := cfg.db.GetUserById(id)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	//TODO: 4) Skicka tillbaka allt utom passwword
	type responseUser struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}
	responseVal := responseUser{
		Id:    user.Id,
		Email: user.Email,
	}
	sendJsonResponse(w, http.StatusOK, responseVal)
}

// POST /api/users
//
//	{
//	    "password": "abc",
//	    "email": "nisse@abc.de"
//	}
func (fdb *apiConfig) handlePostUsers(w http.ResponseWriter, r *http.Request) {
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		sendErrorResponse(w, http.StatusNotAcceptable, err.Error())
		return
	}

	user, err := fdb.db.CreateUser(params.Email, string(hashedPassword))
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	type ResponseUser struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}
	responseVal := ResponseUser{
		Id:    user.Id,
		Email: user.Email,
	}
	sendJsonResponse(w, http.StatusCreated, responseVal)
}

//NOTE: ----------------------------- Authentication ----------------------

// POST /api/login
//
//	{
//	    "password": "abc",
//	    "email": "nisse@abc.de"
//	}
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
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(expires)).UTC()),
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

// PUT /api/users
func (fdb *apiConfig) handleChangeUser(w http.ResponseWriter, r *http.Request) {
	bearer := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	fmt.Println(bearer)
	// TODO:
	// Next, use the jwt.ParseWithClaims function to validate the signature of the JWT
	// and extract the claims into a *jwt.Token struct.
	token, err := jwt.ParseWithClaims(bearer, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return fdb.jwtSecret, nil
	})
	if err != nil {
		// TODO:
		// An error will be returned if the token is invalid or has expired.
		// If the token is invalid, return a 401 Unauthorized response.
		sendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	if !token.Valid {
		sendErrorResponse(w, http.StatusUnauthorized, "Token is not valid")
	}

	// TODO:
	// If all is well with the token,
	// use the token.Claims interface to get access to the user's id from the claims
	// (which should be stored in the Subject field).
	claims := token.Claims.(*jwt.RegisteredClaims)
	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, "ID not found or invalid")
		return
	}
	fmt.Println("ID:", id)

	type ResponseUser struct {
		Id int `json:"id"`
	}
	responseVal := ResponseUser{
		Id: id,
	}
	fmt.Println("ID: ", id)
	sendJsonResponse(w, http.StatusOK, responseVal)
}
