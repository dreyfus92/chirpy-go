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

type returnVals struct{
	Id string `json:"id"`
	Body string `json:"body"`
}

var profanities = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
}


func handlerValidateChirp(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	chirp := paramsValidateChirp{}
	err := decoder.Decode(&chirp)

	const maxChirpLength = 140
	const minChirpLength = 1

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if len(chirp.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest,"Chirp is too long")
		return
	}

	cleaned_body := cleanProfanities(chirp.Body)

	respondWithJSON(w, http.StatusOK, returnVals{
		Id: "1234",
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

