package main

import database "github.com/Schinkenkoenig/chirpy/internal/database"

type apiConfig struct {
	db             *database.DB
	fileserverHits int
}

func (ac *apiConfig) incHits() {
	ac.fileserverHits++
}
