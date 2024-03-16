package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	database "github.com/Schinkenkoenig/chirpy/internal/database"
	"github.com/joho/godotenv"
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
	mux.HandleFunc("PUT /api/users", ac.UpdateUserHandler)
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

	expiresIn := time.Hour * 24

	if request.ExpiresIn != nil {
		fmt.Printf("expires was provided: %d", *request.ExpiresIn)
		seconds := time.Second * time.Duration(*request.ExpiresIn)

		if seconds < expiresIn {
			expiresIn = seconds
		}

		fmt.Printf("expires ultimately is %v", expiresIn)

	}

	tok, err := ac.createJwt(user.Id, expiresIn)
	if err != nil {
		respondWithError(httpWriter, 500, "could not create token")
		fmt.Println(err)
		return
	}

	resp := LoginResponse{Id: user.Id, Email: user.Email, Token: tok}

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

	response := ChirpResponse{
		Id:   chirp.Id,
		Body: chirp.Body,
	}

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

	godotenv.Load()

	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		panic("environment variable JWT_SECRET needs to be set.")
	}

	db, err := database.NewDb(filename)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	ac := apiConfig{fileserverHits: 0, db: db, jwtSecret: jwtSecret}

	Routes(mux, &ac)

	corsMux := middlewareCors(middlewareLog(mux))

	err = http.ListenAndServe(":8080", corsMux)
	if err != nil {
		panic(err)
	}
}
