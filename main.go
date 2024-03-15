package main

import (
	"encoding/json"
	"net/http"
	"strings"

	database "github.com/Schinkenkoenig/chirpy/database"
)

func FileServerHandler(ac *apiConfig) http.Handler {
	httpHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	metricsHttp := ac.middlewareMetricsInc(httpHandler)

	return metricsHttp
}

func Routes(mux *http.ServeMux, ac *apiConfig) {
	mux.Handle("/app/", FileServerHandler(ac))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", ac.handlerMetrics)
	mux.HandleFunc("GET /api/reset", ac.handlerResetMetrics)
	mux.HandleFunc("POST /api/validate_chirp", ValidateChirpHandler)
}

func (db *database.DB) AddChirpHandler(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	// unmarshall
	type Request struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(httpRequest.Body)
	var request Request
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(request.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	request.Body = cleanupChirp(request.Body)

	db.

		// add to db

		resp := Response{CleanedBody: cleanupChirp(_r.Body)}
	respondWithJSON(w, 200, resp)
}

func ValidateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Body string `json:"body"`
	}

	type Response struct {
		CleanedBody string `json:"cleaned_body"`
	}

	d := json.NewDecoder(r.Body)
	var _r Request
	err := d.Decode(&_r)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(_r.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	resp := Response{CleanedBody: cleanupChirp(_r.Body)}
	respondWithJSON(w, 200, resp)
}

func cleanupChirp(s string) string {
	s = strings.ReplaceAll(s, "kerfuffle", "****")
	s = strings.ReplaceAll(s, "sharbert", "****")
	s = strings.ReplaceAll(s, "fornax", "****")
	s = strings.ReplaceAll(s, "Kerfuffle", "****")
	s = strings.ReplaceAll(s, "Sharbert", "****")
	s = strings.ReplaceAll(s, "Fornax", "****")
	return s
}

func main() {
	ac := apiConfig{fileserverHits: 0}
	mux := http.NewServeMux()

	Routes(mux, &ac)

	corsMux := middlewareCors(middlewareLog(mux))

	err := http.ListenAndServe(":8080", corsMux)
	if err != nil {
		panic("oh no")
	}
}
