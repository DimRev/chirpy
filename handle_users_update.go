package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		log.Println("malformed authorization header")
		respondWithError(w, 401, "malformed authorization header")
		return
	}

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
		respondWithError(w, 500, "Error updating user in")
		return
	}

	updatedUser, err := cfg.db.UpdateUser(params.Email, params.Password, splitAuth[1])
	if err != nil {
		log.Printf("Error updating user: %v", err)
		respondWithError(w, 401, "Error updating user in")
		return
	}

	respondWithJSON(w, 200, updatedUser)
}
