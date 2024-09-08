package main

import "net/http"
import "encoding/json"
import "strings"

func (fdb *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := fdb.db.GetChirps()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Could not read chirps")
		return
	}
	sendJsonResponse(w, http.StatusOK, chirps)
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
