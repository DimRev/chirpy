package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/DimRev/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding params: %v", err)
		respondWithError(w, 500, "Error creating user")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Couldn't hash password")
		return
	}

	createdUser, err := cfg.db.CreateUser(params.Email, hashedPassword)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		respondWithError(w, 500, "Error creating user")
		return
	}

	respondWithJSON(w, 201, User{
		ID:          createdUser.Id,
		Email:       createdUser.Email,
		IsChirpyRed: createdUser.IsChirpyRed,
	})
}
