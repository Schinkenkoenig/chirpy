package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func cleanupChirp(s string) string {
	s = strings.ReplaceAll(s, "kerfuffle", "****")
	s = strings.ReplaceAll(s, "sharbert", "****")
	s = strings.ReplaceAll(s, "fornax", "****")
	s = strings.ReplaceAll(s, "Kerfuffle", "****")
	s = strings.ReplaceAll(s, "Sharbert", "****")
	s = strings.ReplaceAll(s, "Fornax", "****")
	return s
}

func (ac *apiConfig) AddChirpHandler(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	// unmarshall
	db := ac.db

	authHeader := httpRequest.Header.Get("Authorization")
	tok := strings.Replace(authHeader, "Bearer ", "", 1)

	userId, err := ac.validateJwt(tok)
	if err != nil {
		respondWithError(httpWriter, 401, "token invalid")
		fmt.Printf("token: '%s' could not be parsded, err: %v\n", tok, err)
		return
	}

	decoder := json.NewDecoder(httpRequest.Body)
	var request chirpRequest
	err = decoder.Decode(&request)
	if err != nil {
		respondWithError(httpWriter, 500, "Something went wrong")
		return
	}

	if len(request.Body) > 140 {
		respondWithError(httpWriter, 400, "Chirp is too long")
		return
	}
	request.Body = cleanupChirp(request.Body)

	chirp, err := db.CreateChirp(userId, request.Body)
	if err != nil {
		respondWithError(httpWriter, 500, "Internal sever error.")
	}

	response := ChirpResponse{
		Id:       chirp.Id,
		Body:     chirp.Body,
		AuthorId: chirp.AuthorId,
	}

	respondWithJSON(httpWriter, 201, response)
}

func (ac *apiConfig) GetChirpByIdHandler(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	// unmarshall
	db := ac.db

	id, err := strconv.Atoi(httpRequest.PathValue("id"))
	if err != nil {
		respondWithError(httpWriter, 400, "Please provide an integer as the id")
		return
	}

	chirp, err := db.GetChirpById(id)
	if err != nil {
		respondWithError(httpWriter, 404, "Could not find chirp.")
		return
	}

	resp := ChirpResponse{Id: chirp.Id, Body: chirp.Body, AuthorId: chirp.AuthorId}

	respondWithJSON(httpWriter, 200, resp)
}

func (ac *apiConfig) GetChirpsHandler(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	// unmarshall
	db := ac.db

	chirps, err := db.GetChirps()
	if err != nil {
		respondWithError(httpWriter, 500, "Internal sever error.")
		fmt.Println(err)
		return
	}

	chirpsResponse := make([]ChirpResponse, 0, len(chirps))

	for _, c := range chirps {
		chirpsResponse = append(chirpsResponse, ChirpResponse{Id: c.Id, Body: c.Body, AuthorId: c.AuthorId})
	}

	respondWithJSON(httpWriter, 200, chirpsResponse)
}
