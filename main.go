package main

import (
	"encoding/json"
	"log"
	"net/http"

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

	createdUser, err := cfg.db.Login(params.Email, params.Password)
	if err != nil {
		log.Printf("Error logging in user: %v", err)
		respondWithError(w, 401, "Error logging in")
		return
	}

	respondWithJSON(w, 200, createdUser)
}
