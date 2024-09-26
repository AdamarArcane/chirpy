package main

import (
	"net/http"
	"sync/atomic"
)

func main() {
	mux := http.NewServeMux()

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	fileServerHandler := http.FileServer(http.Dir("."))
	strippedHandler := http.StripPrefix("/app", fileServerHandler)

	mux.Handle("/app/", cfg.siteHitsMiddleware(strippedHandler))
	mux.Handle("/app/assets/logo.png", strippedHandler)
	mux.HandleFunc("GET /api/healthz", handlerHealthz)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	server.ListenAndServe()
}

// Data Structures

type apiConfig struct {
	fileserverHits atomic.Int32
}

// Handler Functions
