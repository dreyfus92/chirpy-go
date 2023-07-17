package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type paramsValidateChirp struct {
	Body string `json:"body"`
}

type returnVals struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

var profanities = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
}

// This handler validates creates and validates chirps
func (apiCfg apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	chirp := paramsValidateChirp{}
	err := decoder.Decode(&chirp)
	const maxChirpLength = 140

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if len(chirp.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	// Clean profanities
	cleaned_body := cleanProfanities(chirp.Body)

	chirp_id := apiCfg.db.CreateChirp(cleaned_body)

	respondWithJSON(w, http.StatusCreated, returnVals{
		Id:   chirp_id,
		Body: cleaned_body,
	})
}

func cleanProfanities(body string) string {
	bodyWords := strings.Split(body, " ")
	log.Printf("Body words: %v", bodyWords)
	for i, word := range bodyWords {
		for _, profanity := range profanities {
			if strings.Contains(strings.ToLower(word), profanity) {
				bodyWords[i] = strings.Repeat("*", 4)
			}
		}
	}
	return strings.Join(bodyWords, " ")
}
