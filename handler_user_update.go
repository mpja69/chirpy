package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// NOTE: handleUpdateUser, PUT /api/users, authorizes and updates a users info
func (fdb *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	bearer := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

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
