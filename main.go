package main

import (
	"flag"
	"net/http"
	"os"

	database "github.com/Schinkenkoenig/chirpy/internal/database"
	"github.com/joho/godotenv"
)

func FileServerHandler(ac *apiConfig) http.Handler {
	httpHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	metricsHttp := ac.middlewareMetricsInc(httpHandler)

	return metricsHttp
}

func Routes(mux *http.ServeMux, ac *apiConfig) {
	// app
	mux.Handle("/app/", FileServerHandler(ac))

	// admin
	mux.HandleFunc("GET /admin/metrics", ac.handlerMetrics)

	// api-metrics
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", ac.handlerResetMetrics)

	// api/chirps
	mux.HandleFunc("POST /api/chirps", ac.AddChirpHandler)
	mux.HandleFunc("GET /api/chirps", ac.GetChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{id}", ac.GetChirpByIdHandler)

	// api/login
	mux.HandleFunc("POST /api/login", ac.LoginUserHandler)
	mux.HandleFunc("POST /api/refresh", ac.RefreshTokenHandler)
	mux.HandleFunc("POST /api/revoke", ac.RevokeTokenHandler)

	// api/users
	mux.HandleFunc("POST /api/users", ac.AddUserHandler)
	mux.HandleFunc("PUT /api/users", ac.UpdateUserHandler)
}

func main() {
	filename := "database.json"

	// debug flag parsing
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *dbg {
		os.Remove(filename)
	}

	// env parsing
	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("environment variable JWT_SECRET needs to be set.")
	}

	// db creation
	db, err := database.NewDb(filename)
	if err != nil {
		panic(err)
	}

	ac := apiConfig{fileserverHits: 0, db: db, jwtSecret: jwtSecret}

	// http mux
	mux := http.NewServeMux()
	Routes(mux, &ac)
	corsMux := middlewareCors(middlewareLog(mux))

	err = http.ListenAndServe(":8080", corsMux)
	if err != nil {
		panic(err)
	}
}
