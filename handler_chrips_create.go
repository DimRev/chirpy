package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {
	bannedWords := []string{
		"Kerfuffle",
		"kerfuffle",
		"Sharbert",
		"sharbert",
		"Fornax",
		"fornax",
	}

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(params.Body) > 140 {
		w.WriteHeader(400)
		return
	}

	formattedWords := []string{}

	words := strings.Fields(params.Body)
	for _, word := range words {
		if slices.Contains(bannedWords, word) {
			formattedWords = append(formattedWords, "****")
		} else {
			formattedWords = append(formattedWords, word)
		}
	}

	formattedBody := strings.Join(formattedWords, " ")

	createdChirp, err := cfg.db.CreateChirp(formattedBody)
	if err != nil {
		log.Printf("Error creating chirp: %v", err)
		w.WriteHeader(500)
		return
	}

	respondWithJSON(w, 201, createdChirp)
}
