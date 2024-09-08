package main

import "encoding/json"
import "log"
import "net/http"

func sendErrorResponse(w http.ResponseWriter, code int, msg string) {
	log.Printf("Error: %d", code)
	type returnVals struct {
		Error string `json:"error"`
	}
	sendJsonResponse(w, code, returnVals{
		Error: msg,
	})
}

func sendJsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
