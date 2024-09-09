package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"github.com/mpja69/chirpy/internal/database"
	"golang.org/x/crypto/bcrypt"
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
	//TODO: 1) Ta in ett password från "payload:en"
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

	//TODO: 2) Kryptera password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		sendErrorResponse(w, http.StatusNotAcceptable, err.Error())
		return
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(params.Password))
	if err != nil {
		sendErrorResponse(w, http.StatusNotAcceptable, err.Error())
	}

	//TODO: 3) Spara det krypterade pw i db
	user, err := fdb.db.CreateUser(params.Email, string(hashedPassword))
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
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
	sendJsonResponse(w, http.StatusCreated, responseVal)
}
