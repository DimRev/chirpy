package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	authorIdStr := r.URL.Query().Get("author_id")
	chirps, err := cfg.db.GetChirps(authorIdStr)
	if err != nil {
		log.Printf("Error getting chirp: %v", err)
		respondWithError(w, 500, "Error getting chirp")
		return
	}

	respondWithJSON(w, 200, chirps)
}
