package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding params: %v", err)
		respondWithError(w, 500, "Error logging in")
		return
	}

	loggedInUser, err := cfg.db.Login(params.Email, params.Password, params.ExpiresInSeconds)
	if err != nil {
		log.Printf("Error logging in user: %v", err)
		respondWithError(w, 401, "Error logging in")
		return
	}

	respondWithJSON(w, 200, loggedInUser)
}
