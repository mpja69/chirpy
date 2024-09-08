package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/mpja69/chirpy/internal/database"
)

func (fdb *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := fdb.db.GetChirps()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Could not read chirps")
		return
	}
	sortedChirps := []database.Chirp{}
	for _, ch := range chirps {
		sortedChirps = append(sortedChirps, database.Chirp{
			Id:   ch.Id,
			Body: ch.Body,
		})
	}
	sort.Slice(sortedChirps, func(i int, j int) bool {
		return sortedChirps[i].Id < sortedChirps[j].Id
	})

	sendJsonResponse(w, http.StatusOK, sortedChirps)
}

func (cfg *apiConfig) handleGetChirpById(w http.ResponseWriter, r *http.Request) {
	idValue := r.PathValue("chirpId")
	id, err := strconv.Atoi(idValue)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "The requested id is malformed")
		return
	}
	chirp, err := cfg.db.GetChirp(id)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, "Chirp does not exist")
		return
	}
	sendJsonResponse(w, http.StatusOK, chirp)
}

func (fdb *apiConfig) handlePostChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Something went wrong")
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

	chirp, err := fdb.db.CreateChirp(cleanedMsg)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Could not create chirp")
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
