package main

type apiConfig struct {
	fileserverHits int
}

func (ac *apiConfig) incHits() {
	ac.fileserverHits++
}
