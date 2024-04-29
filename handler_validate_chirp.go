package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
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

	formattedString := strings.Join(formattedWords, " ")

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	respBody := returnVals{
		CleanedBody: formattedString,
	}

	respondWithJSON(w, 200, respBody)
}
