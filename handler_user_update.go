package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mpja69/chirpy/internal/auth"

	"golang.org/x/crypto/bcrypt"
)

// NOTE: handleUpdateUser, PUT /api/users, authorizes and updates a users info
func (fdb *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	accessToken, err := auth.GetBearerToken(r)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	userIdString, err := auth.ValidateJWT(accessToken, fdb.jwtSecret)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, err.Error())
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

	user, err := fdb.db.GetUserById(id)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
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
	sendJsonResponse(w, http.StatusOK, responseVal)
}

// handleWebhooks - "POST /api/polka/webhooks"
func (cfg *apiConfig) handleWebhooks(w http.ResponseWriter, r *http.Request) {
	type userData struct {
		UserId int `json:"user_id"`
	}
	type parameters struct {
		Event string   `json:"event"`
		Data  userData `json:"data"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	_, err = cfg.db.UpgradeUser(params.Data.UserId, true)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	sendJsonResponse(w, http.StatusNoContent, struct{}{})

}
