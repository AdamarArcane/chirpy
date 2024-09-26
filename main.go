package main

import (
	"fmt"
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
	mux.HandleFunc("/healthz", handlerHealthz)
	mux.HandleFunc("/metrics", cfg.handlerHits)
	mux.HandleFunc("/reset", cfg.handlerReset)

	server.ListenAndServe()
}

// Data Structures

type apiConfig struct {
	fileserverHits atomic.Int32
}

// Handler Functions

func handlerHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// Middleware

func (cfg *apiConfig) siteHitsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	response := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	w.Write([]byte(response))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = atomic.Int32{}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	response := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	w.Write([]byte(response))
}
