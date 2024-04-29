package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps()
	if err != nil {
		log.Printf("Error getting chirp: %v", err)
		respondWithError(w, 500, "Error getting chirp")
		return
	}

	respondWithJSON(w, 200, chirps)
}
