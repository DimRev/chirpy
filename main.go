package main

import (
	"log"
	"net/http"
)

const (
	PORT           = "8080"
	FILE_PATH_ROOT = "."
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)

	cfg := apiConfig{
		fileserverHits: 0,
	}

	srv := &http.Server{
		Addr:    ":" + PORT,
		Handler: corsMux,
	}

	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("/healthz", handlerHealthCheck)
	mux.HandleFunc("/metrics", cfg.handlerMetrics)
	mux.HandleFunc("/reset", cfg.handlerReset)

	log.Printf("Serving files from %s on http://localhost:%s\n", FILE_PATH_ROOT, PORT)
	log.Fatal(srv.ListenAndServe())
}
