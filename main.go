package main

import (
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

	log.Printf("Serving files from %s on http://localhost:%s\n", FILE_PATH_ROOT, PORT)
	log.Fatal(srv.ListenAndServe())
}
