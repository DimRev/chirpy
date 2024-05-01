package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/DimRev/chirpy/internal/auth"
)

func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		log.Println("malformed authorization header")
		respondWithError(w, 401, "malformed authorization header")
		return
	}
	tokenString := splitAuth[1]

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	fmt.Println()

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding params: %v", err)
		respondWithError(w, 500, "Couldn't decode request body")
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

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Couldn't hash password")
		return
	}

	updatedUser, err := cfg.db.UpdateUser(params.Email, hashedPassword, userIDInt)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		respondWithError(w, 401, "Error updating user in")
		return
	}

	respondWithJSON(w, 200, updatedUser)
}
