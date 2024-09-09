package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

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

// POST /api/login
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
	}

	type ResponseUser struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}
	responseVal := ResponseUser{
		Id:    user.Id,
		Email: user.Email,
	}
	sendJsonResponse(w, http.StatusOK, responseVal)
}
