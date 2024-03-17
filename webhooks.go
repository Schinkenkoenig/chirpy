package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (ac *apiConfig) webhookHandler(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	// unmarshall
	db := ac.db

	decoder := json.NewDecoder(httpRequest.Body)
	var request webhookRequest
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(httpWriter, 500, "Something went wrong")
		return
	}

	authHeader := httpRequest.Header.Get("Authorization")
	apiKey := strings.Replace(authHeader, "ApiKey ", "", 1)

	if apiKey != ac.polkaKey {
		respondWithError(httpWriter, 401, "not authenticated")
		return
	}

	if request.Event != "user.upgraded" {
		httpWriter.WriteHeader(200)
		httpWriter.Write([]byte{})
		return
	}

	_, err = db.UpgradeUser(request.Data.UserId)
	if err != nil {
		httpWriter.WriteHeader(404)
		httpWriter.Write([]byte{})
		return
	}

	httpWriter.WriteHeader(200)
	httpWriter.Write([]byte{})
}
