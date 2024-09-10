package main

import (
	"encoding/json"
	"fmt"
	"github.com/mpja69/chirpy/internal/auth"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// NOTE: handleUpdateUser, PUT /api/users, authorizes and updates a users info
func (fdb *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	token, err := auth.GetBearerToken(r, string(fdb.jwtSecret))
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	userIdString, err := auth.ValidateJWT(token, fdb.jwtSecret)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := strconv.Atoi(userIdString)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, "ID not a number")
		return
	}
	fmt.Println("ID:", id)

	//----------------------- Update the User with the data in body -----------------------------
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		sendErrorResponse(w, http.StatusNotAcceptable, err.Error())
		return
	}

	_, err = fdb.db.UpdateUser(id, params.Email, string(hashedPassword))
	if err != nil {
		sendErrorResponse(w, http.StatusNotAcceptable, err.Error())
		return
	}

	type ResponseUser struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}
	responseVal := ResponseUser{
		Id:    id,
		Email: params.Email,
	}
	fmt.Println("ID: ", id)
	sendJsonResponse(w, http.StatusOK, responseVal)
}
