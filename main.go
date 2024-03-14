package main

import (
	"net/http"
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
