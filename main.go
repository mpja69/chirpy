package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fileRoot := "."
	port := "8080"
	mux := http.NewServeMux()
	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	// Add a handler for files, starting in root
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir(fileRoot))))

	// Add a handleFunc for a specific path
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))

	})

	fmt.Printf("Starting server on: %s, serving files from: %s\n", srv.Addr, fileRoot)
	log.Fatal(srv.ListenAndServe())

}
