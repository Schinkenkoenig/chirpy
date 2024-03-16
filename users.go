package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (ac *apiConfig) UpdateUserHandler(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	// unmarshall
	db := ac.db

	decoder := json.NewDecoder(httpRequest.Body)
	var request userRequest
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(httpWriter, 500, "Something went wrong")
		return
	}

	authHeader := httpRequest.Header.Get("Authorization")
	tok := strings.Replace(authHeader, "Bearer ", "", 1)

	userId, err := ac.validateJwt(tok)
	if err != nil {
		respondWithError(httpWriter, 401, "token invalid")
		fmt.Printf("token: '%s' could not be parsded, err: %v\n", tok, err)
		return
	}

	updated, err := db.UpdateUser(userId, request.Email, request.Password)
	if err != nil {
		respondWithError(httpWriter, 404, "not found")
		return
	}

	resp := UserResponse{Email: updated.Email, Id: updated.Id}

	respondWithJSON(httpWriter, 200, resp)
}

func (ac *apiConfig) AddUserHandler(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	// unmarshall
	db := ac.db

	decoder := json.NewDecoder(httpRequest.Body)
	var request userRequest
	err := decoder.Decode(&request)
	if err != nil {
		fmt.Printf("%v", err)
		respondWithError(httpWriter, 500, "Something went wrong")
		return
	}

	user, err := db.CreateUser(request.Email, request.Password)
	if err != nil {
		fmt.Printf("%v", err)
		respondWithError(httpWriter, 500, "Internal sever error.")
		return
	}

	resp := UserResponse{Id: user.Id, Email: user.Email}

	respondWithJSON(httpWriter, 201, resp)
}
