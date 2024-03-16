package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type RefreshTokenResponse struct {
	Token string `json:"token"`
}

func (ac *apiConfig) RefreshTokenHandler(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	// unmarshall
	db := ac.db

	authHeader := httpRequest.Header.Get("Authorization")
	tok := strings.Replace(authHeader, "Bearer ", "", 1)

	userId, err := ac.validateRefreshJwt(tok)
	if err != nil {
		respondWithError(httpWriter, 401, "token invalid")
		fmt.Printf("token: '%s' could not be parsded, err: %v\n", tok, err)
		return
	}

	// todo: check revocations
	err = db.IsTokenRevoked(userId, tok)
	if err != nil {
		respondWithError(httpWriter, 401, "revoked")
	}

	jwt, err := ac.createJwt(userId)
	if err != nil {
		respondWithError(httpWriter, 500, "could not create token")
		fmt.Println(err)
		return
	}

	resp := RefreshTokenResponse{Token: jwt.accessToken}

	respondWithJSON(httpWriter, 200, resp)
}

func (ac *apiConfig) RevokeTokenHandler(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	// unmarshall
	db := ac.db

	authHeader := httpRequest.Header.Get("Authorization")
	tok := strings.Replace(authHeader, "Bearer ", "", 1)

	userId, err := ac.validateRefreshJwt(tok)
	if err != nil {
		respondWithError(httpWriter, 401, "token invalid")
		fmt.Printf("token: '%s' could not be parsded, err: %v\n", tok, err)
		return
	}

	db.RevokeToken(userId, tok)
	var resp string

	respondWithJSON(httpWriter, 200, resp)
}

func (ac *apiConfig) LoginUserHandler(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	// unmarshall
	db := ac.db

	decoder := json.NewDecoder(httpRequest.Body)
	var request loginRequest
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(httpWriter, 500, "Something went wrong")
		return
	}

	user, err := db.IsPasswordCorrect(request.Email, request.Password)
	if err != nil {
		respondWithError(httpWriter, 401, "Internal sever error.")
		return
	}

	tok, err := ac.createJwt(user.Id)
	if err != nil {
		respondWithError(httpWriter, 500, "could not create token")
		fmt.Println(err)
		return
	}

	resp := LoginResponse{Id: user.Id, Email: user.Email, AccessToken: tok.accessToken, RefreshToken: tok.refreshToken}

	respondWithJSON(httpWriter, 200, resp)
}
