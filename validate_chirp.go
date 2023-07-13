package main

import (
	"encoding/json"
	"net/http"
)

type paramsValidateChirp struct {
	Body string `json:"body"`
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	chirp := paramsValidateChirp{}

	err := decoder.Decode(&chirp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Something went wrong"}`))
		return
	}

	if len(chirp.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Chirp is too long"}`))
		return
	}


	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"valid": true}`))
}
