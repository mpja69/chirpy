package main

import (
	"fmt"
	"log"
	"net/http"
)

type fileHandler struct {
}

func (h *fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func main() {
	port := "8080"
	mux := http.NewServeMux()
	srv := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	mux.Handle("/", http.FileServer(http.Dir(".")))

	fmt.Printf("Starting server on %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())

}
