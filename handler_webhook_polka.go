package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		log.Println("malformed authorization header")
		respondWithError(w, 401, "malformed authorization header")
		return
	}
	apiKey := splitAuth[1]
	if apiKey != cfg.polkaApiKey {
		log.Println("invalid apiKey")
		respondWithError(w, 401, "invalid apiKey")
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId int `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding params: %v", err)
		respondWithError(w, 500, "Failed decoding request body")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, 200, "")
		return
	}

	upgradedUser, err := cfg.db.UpdateChirpyRed(params.Data.UserId)
	if err != nil {
		log.Printf("Error upgrading to chirpy red: %v", err)
		respondWithError(w, 500, "Failed to upgrade to chirpy red")
		return
	}

	respondWithJSON(w, 200, User{
		ID:          upgradedUser.Id,
		Email:       upgradedUser.Email,
		IsChirpyRed: upgradedUser.IsChirpyRed,
	})
}
