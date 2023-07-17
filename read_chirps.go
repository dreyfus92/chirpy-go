package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (apiCfg apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	allChirps, err := apiCfg.db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps")
		log.Fatal(err)
	}
	respondWithJSON(w, http.StatusOK, allChirps)
}

func (apiCfg apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	// get the ID from the URL
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	// find the chirp from the id if possible
	chirp, err := apiCfg.db.GetChirp(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	// respond with the chirp
	respondWithJSON(w, http.StatusOK, chirp)
}
