package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func (apiCfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	// decode the user from JSON into struct
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// create the user
	email := params.Email
	newUser := apiCfg.db.CreateUser(email)
	fmt.Printf("Created user: %v", newUser)

	// respond with the user
	respondWithJSON(w, 201, User{
		Id:    newUser,
		Email: email,
	})

}
