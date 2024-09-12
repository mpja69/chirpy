package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// handleGetUsers - GET /api/users
func (fdb *apiConfig) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := fdb.db.GetUsers()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	type ResponseUser struct {
		Id          int    `json:"id"`
		Email       string `json:"email"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}
	sortedUsers := []ResponseUser{}
	for _, u := range users {
		sortedUsers = append(sortedUsers, ResponseUser{
			Id:          u.Id,
			Email:       u.Email,
			IsChirpyRed: u.IsChirpyRed,
		})
	}
	sort.Slice(sortedUsers, func(i int, j int) bool {
		return sortedUsers[i].Id < sortedUsers[j].Id
	})

	sendJsonResponse(w, http.StatusOK, sortedUsers)
}

// handleGetUserById -  GET /api/users/{userId}
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
		Id          int    `json:"id"`
		Email       string `json:"email"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}
	responseVal := responseUser{
		Id:          user.Id,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	sendJsonResponse(w, http.StatusOK, responseVal)
}

// handleCreateUser - POST /api/users
func (fdb *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
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
		Id          int    `json:"id"`
		Email       string `json:"email"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}
	responseVal := ResponseUser{
		Id:          user.Id,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	sendJsonResponse(w, http.StatusCreated, responseVal)
}
