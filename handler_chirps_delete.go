package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/DimRev/chirpy/internal/auth"
)

func (cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		log.Println("malformed authorization header")
		respondWithError(w, 401, "malformed authorization header")
		return
	}
	tokenString := splitAuth[1]

	chirpIdStr := r.PathValue("id")

	chirpIdInt, err := strconv.Atoi(chirpIdStr)
	if err != nil {
		log.Printf("Not valid id: %v", err)
		respondWithError(w, 500, "Error getting chirp")
		return
	}

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

	err = cfg.db.DeleteChirp(chirpIdInt, userIDInt)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		respondWithError(w, 403, "Error deleting user")
		return
	}

	respondWithJSON(w, 200, struct{}{})
}
