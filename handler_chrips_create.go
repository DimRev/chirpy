package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/DimRev/chirpy/internal/auth"
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
		log.Printf("Error decoding params: %v", err)
		respondWithError(w, 500, "Error creating chirp")
		return
	}

	authHeader := r.Header.Get("Authorization")
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		log.Println("malformed authorization header")
		respondWithError(w, 401, "malformed authorization header")
		return
	}
	tokenString := splitAuth[1]

	userIdString, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error decoding params: %v", err)
		respondWithError(w, 401, "Couldn't validate token")
		return
	}

	userIDInt, err := strconv.Atoi(userIdString)
	if err != nil {
		respondWithError(w, 500, "Couldn't parse user ID")
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

	createdChirp, err := cfg.db.CreateChirp(formattedBody, userIDInt)
	if err != nil {
		log.Printf("Error creating chirp: %v", err)
		respondWithError(w, 500, "Error creating chirp")
		return
	}

	respondWithJSON(w, 201, createdChirp)
}
