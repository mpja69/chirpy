package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/mpja69/chirpy/internal/auth"
	"github.com/mpja69/chirpy/internal/database"
)

// handleGetChirps "GET /api/chirps",
func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sortedChirps := []database.Chirp{}
	for _, chirp := range chirps {
		sortedChirps = append(sortedChirps, database.Chirp{
			Id:       chirp.Id,
			Body:     chirp.Body,
			AuthorId: chirp.AuthorId,
		})
	}
	sort.Slice(sortedChirps, func(i int, j int) bool {
		return sortedChirps[i].Id < sortedChirps[j].Id
	})

	sendJsonResponse(w, http.StatusOK, sortedChirps)
}

// handleGetChirpById - "GET /api/chirps/{chirpId}"
func (cfg *apiConfig) handleGetChirpById(w http.ResponseWriter, r *http.Request) {
	idValue := r.PathValue("chirpId")
	id, err := strconv.Atoi(idValue)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	chirp, err := cfg.db.GetChirp(id)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	sendJsonResponse(w, http.StatusOK, chirp)
}

// handlePostChirps - "POST /api/chirps"
func (cfg *apiConfig) handlePostChirps(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetAuthorizationBearer(r)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	userIdString, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "ID not a number")
		return
	}
	user, err := cfg.db.GetUserById(userId)
	log.Printf("HandlePostChirps: User: %d, %s is authenticated\n", user.Id, user.Email)

	type parameters struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(params.Body) > 140 {
		sendErrorResponse(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleanedMsg := cleanBody(params.Body, profaneWords)

	chirp, err := cfg.db.CreateChirp(cleanedMsg, user.Id)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendJsonResponse(w, http.StatusCreated, chirp)
}

func cleanBody(msg string, badWords map[string]struct{}) string {
	words := strings.Split(msg, " ")
	for i, word := range words {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}
	cleanedMsg := strings.Join(words, " ")
	return cleanedMsg
}

// handleDeleteChirps - "DELETE /api/chirps/{chirpID}"
func (cfg *apiConfig) handleDeleteChirpById(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetAuthorizationBearer(r)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	userIdString, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "ID not a number")
		return
	}
	user, err := cfg.db.GetUserById(userId)
	log.Printf("HandleDeleteChirps: User: %d, %s is authenticated\n", user.Id, user.Email)

	idValue := r.PathValue("chirpId")
	id, err := strconv.Atoi(idValue)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	chirp, err := cfg.db.GetChirp(id)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	if chirp.AuthorId != user.Id {
		w.WriteHeader(http.StatusForbidden)
	}
	sendJsonResponse(w, http.StatusNoContent, chirp)
}
