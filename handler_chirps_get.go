package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps()
	if err != nil {
		log.Printf("Error creating chirp: %v", err)
		w.WriteHeader(500)
		return
	}

	respondWithJSON(w, 200, chirps)
}
