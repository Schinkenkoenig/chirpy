package main

import (
	"fmt"
	"net/http"
)

func (ac *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ac.incHits()
		next.ServeHTTP(w, r)
	})
}

func (ac *apiConfig) handlerResetMetrics(w http.ResponseWriter, _ *http.Request) {
	ac.fileserverHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	var body []byte
	w.Write(body)
}

func (ac *apiConfig) handlerMetrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	bodyString := fmt.Sprintf(`<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>`, ac.fileserverHits)
	body := []byte(bodyString)
	w.Write(body)
}
