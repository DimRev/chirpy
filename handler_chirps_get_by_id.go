package main

import (
	"log"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) handleGetChirpById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Not valid id: %v", err)
		respondWithError(w, 500, "Error getting chirp")
		return
	}

	chirp, err := cfg.db.GetChirpById(id)
	if err != nil {
		log.Printf("Problem getting chirp from DB: %v", err)
		respondWithError(w, 404, "Not found")
		return
	}

	respondWithJSON(w, 200, chirp)
}
