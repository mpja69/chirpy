package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"github.com/mpja69/chirpy/internal/database"
)

func (fdb *apiConfig) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := fdb.db.GetUsers()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sortedUsers := []database.User{}
	for _, u := range users {
		sortedUsers = append(sortedUsers, database.User{
			Id:    u.Id,
			Email: u.Email,
		})
	}
	sort.Slice(sortedUsers, func(i int, j int) bool {
		return sortedUsers[i].Id < sortedUsers[j].Id
	})

	sendJsonResponse(w, http.StatusOK, sortedUsers)
}

func (cfg *apiConfig) handleGetUserById(w http.ResponseWriter, r *http.Request) {
	idValue := r.PathValue("userId")
	id, err := strconv.Atoi(idValue)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := cfg.db.GetUser(id)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	sendJsonResponse(w, http.StatusOK, user)
}

func (fdb *apiConfig) handlePostUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := fdb.db.CreateUser(params.Email)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendJsonResponse(w, http.StatusCreated, user)
}
