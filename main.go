package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits += 1
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func main() {
	fileRoot := "."
	port := "8080"
	// Add a handler for files, starting in root
	// mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir(fileRoot))))
	apiCfg := apiConfig{}
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(fileRoot)))))

	// Add a handleFunc for a specific path
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))

	})

	mux.HandleFunc("/reset", apiCfg.handlerReset)

	mux.HandleFunc("/metrics", apiCfg.handlerMetrics)

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	fmt.Printf("Starting server on: %s, serving files from: %s\n", srv.Addr, fileRoot)
	log.Fatal(srv.ListenAndServe())

}
