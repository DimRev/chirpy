package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
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
