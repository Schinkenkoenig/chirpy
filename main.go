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
	mux.HandleFunc("POST /api/login", ac.LoginUserHandler)
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

	resp := ChirpResponse{Id: chirp.Id, Body: chirp.Body}

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
		chirpsResponse = append(chirpsResponse, ChirpResponse{Id: c.Id, Body: c.Body})
	}

	respondWithJSON(httpWriter, 200, chirpsResponse)
}

func (ac *apiConfig) LoginUserHandler(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	// unmarshall
	db := ac.db

	decoder := json.NewDecoder(httpRequest.Body)
	var request userRequest
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

	resp := UserResponse{Id: user.Id, Email: user.Email}

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

func (ac *apiConfig) AddChirpHandler(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	// unmarshall
	db := ac.db

	decoder := json.NewDecoder(httpRequest.Body)
	var request chirpRequest
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

	response := ChirpResponse{Id: chirp.Id, Body: chirp.Body}

	respondWithJSON(httpWriter, 201, response)
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
