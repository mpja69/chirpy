package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/mpja69/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	db             *database.DB
	jwtSecret      []byte
}

func main() {
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *dbg {
		os.Remove("database.json")
	}

	fileRoot := "."
	port := "8080"
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal("Couldn't create database")
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Env file error: %v\n", err)
	}
	apiCfg := apiConfig{
		fileserverHits: 0,
		db:             db,
		jwtSecret:      []byte(os.Getenv("JWT_SECRET")),
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(fileRoot)))))

	// Add a handleFunc for a specific path
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))

	})
	mux.HandleFunc("/api/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlePostChirps)
	mux.HandleFunc("GET /api/chirps", apiCfg.handleGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpId}", apiCfg.handleGetChirpById)
	mux.HandleFunc("DELETE /api/chirps/{chirpId}", apiCfg.handleDeleteChirpById)

	mux.HandleFunc("POST /api/users", apiCfg.handleCreateUser)
	mux.HandleFunc("GET /api/users", apiCfg.handleGetUsers)
	mux.HandleFunc("GET /api/users/{userId}", apiCfg.handleGetUserById)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handleWebhooks)

	mux.HandleFunc("POST /api/login", apiCfg.handleLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handleRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.handleRevokeToken)
	mux.HandleFunc("PUT /api/users", apiCfg.handleUpdateUser)
	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	fmt.Printf("Starting server on: %s, serving files from: %s\n", srv.Addr, fileRoot)
	log.Fatal(srv.ListenAndServe())

}
