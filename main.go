package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/DimRev/chirpy/internal/database"
)

const (
	PORT           = "8080"
	FILE_PATH_ROOT = "."
)

type apiConfig struct {
	fileserverHits int
	db             *database.DB
}

func main() {
	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)

	db, err := database.NewDB("internal/database/database.json")
	if err != nil {
		log.Printf("Error connecting to DB: %v", err)
	}

	cfg := apiConfig{
		fileserverHits: 0,
		db:             db,
	}

	srv := &http.Server{
		Addr:    ":" + PORT,
		Handler: corsMux,
	}

	mux.Handle("GET /app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /admin/metrics", cfg.handlerAdminMetrics)

	mux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)
	mux.HandleFunc("GET /api/healthz", handlerHealthCheck)
	mux.HandleFunc("GET /api/reset", cfg.handlerMetricsReset)

	mux.HandleFunc("POST /api/chirps", cfg.handleCreateChirp)
	mux.HandleFunc("GET /api/chirps", cfg.handleGetChirps)
	mux.HandleFunc("GET /api/chirps/{id}", cfg.handleGetChirpById)

	mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", cfg.handleUpdateUser)
	mux.HandleFunc("POST /api/login", cfg.handlerLogin)

	log.Printf("Serving files from %s on http://localhost:%s\n", FILE_PATH_ROOT, PORT)
	log.Fatal(srv.ListenAndServe())
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding params: %v", err)
		respondWithError(w, 500, "Error logging in")
		return
	}

	loggedInUser, err := cfg.db.Login(params.Email, params.Password)
	if err != nil {
		log.Printf("Error logging in user: %v", err)
		respondWithError(w, 401, "Error logging in")
		return
	}

	respondWithJSON(w, 200, loggedInUser)
}

func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		log.Println("malformed authorization header")
		respondWithError(w, 401, "malformed authorization header")
		return
	}
	log.Printf(`
********
	Header received: %v
	Fields, 1: %v, 2:%v, len:%v
********`, authHeader, splitAuth[0], splitAuth[1], len(splitAuth))

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	fmt.Println()

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding params: %v", err)
		respondWithError(w, 500, "Error updating user in")
		return
	}

	updatedUser, err := cfg.db.UpdateUser(params.Email, params.Password, splitAuth[1])
	if err != nil {
		log.Printf("Error updating user: %v", err)
		respondWithError(w, 401, "Error updating user in")
		return
	}

	respondWithJSON(w, 200, updatedUser)
}
