package main

import (
	"log"
	"net/http"
	"os"

	"github.com/DimRev/chirpy/internal/database"
	"github.com/joho/godotenv"
)

const (
	PORT           = "8080"
	FILE_PATH_ROOT = "."
)

type apiConfig struct {
	fileserverHits int
	db             *database.DB
	jwtSecret      string
	polkaApiKey    string
}

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"-"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}

func main() {
	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)

	db, err := database.NewDB("internal/database/database.json")
	if err != nil {
		log.Printf("Error connecting to DB: %v", err)
	}

	godotenv.Load(".env")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	polkaApiKey := os.Getenv("POLKA_API_KEY")
	if jwtSecret == "" {
		log.Fatal("POLKA_API_KEY environment variable is not set")
	}

	cfg := apiConfig{
		fileserverHits: 0,
		db:             db,
		jwtSecret:      jwtSecret,
		polkaApiKey:    polkaApiKey,
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
	mux.HandleFunc("DELETE /api/chirps/{id}", cfg.handleDeleteChirp)

	mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", cfg.handleUpdateUser)

	mux.HandleFunc("POST /api/login", cfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevoke)

	mux.HandleFunc("POST /api/polka/webhooks", cfg.handlerPolkaWebhook)

	log.Printf("Serving files from %s on http://localhost:%s\n", FILE_PATH_ROOT, PORT)
	log.Fatal(srv.ListenAndServe())
}
