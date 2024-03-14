package main

import (
	"fmt"
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	body := []byte(http.StatusText(http.StatusOK))
	w.Write(body)

	fmt.Println("Healthz handler")
}
