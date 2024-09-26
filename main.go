package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/adamararcane/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	mux := http.NewServeMux()
	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error opening db: %s", err)
		return
	}

	dbQueries := database.New(db)

	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQ:            dbQueries,
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
	dbQ            *database.Queries
}
