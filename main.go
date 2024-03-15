package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	database "github.com/Schinkenkoenig/chirpy/internal/database"
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
	mux.HandleFunc("POST /api/chirps", ac.AddChirpHandler)
	mux.HandleFunc("POST /api/users", ac.AddUserHandler)
	mux.HandleFunc("GET /api/chirps", ac.GetChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{id}", ac.GetChirpByIdHandler)
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

	respondWithJSON(httpWriter, 200, chirp)
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

	fmt.Printf("%v", chirps)

	respondWithJSON(httpWriter, 200, chirps)
}

func (ac *apiConfig) AddUserHandler(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	// unmarshall
	db := ac.db
	type Request struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(httpRequest.Body)
	var request Request
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(httpWriter, 500, "Something went wrong")
		return
	}

	user, err := db.CreateUser(request.Email)
	if err != nil {
		respondWithError(httpWriter, 500, "Internal sever error.")
	}

	respondWithJSON(httpWriter, 201, user)
}

func (ac *apiConfig) AddChirpHandler(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	// unmarshall
	db := ac.db
	type Request struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(httpRequest.Body)
	var request Request
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(httpWriter, 500, "Something went wrong")
		return
	}

	if len(request.Body) > 140 {
		respondWithError(httpWriter, 400, "Chirp is too long")
		return
	}
	request.Body = cleanupChirp(request.Body)

	chirp, err := db.CreateChirp(request.Body)
	if err != nil {
		respondWithError(httpWriter, 500, "Internal sever error.")
	}

	respondWithJSON(httpWriter, 201, chirp)
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
	filename := "database.json"
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *dbg {
		os.Remove(filename)
	}

	db, err := database.NewDb(filename)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	ac := apiConfig{fileserverHits: 0, db: db}

	Routes(mux, &ac)

	corsMux := middlewareCors(middlewareLog(mux))

	err = http.ListenAndServe(":8080", corsMux)
	if err != nil {
		panic(err)
	}
}
