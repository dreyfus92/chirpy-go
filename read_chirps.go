package main

import (
	"log"
	"net/http"
)

func (apiCfg apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	allChirps, err := apiCfg.db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps")
		log.Fatal(err)
	}
	respondWithJSON(w, http.StatusOK, allChirps)
}
