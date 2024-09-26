package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		WriteErrorResponse(w, 500, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		log.Printf("Chirp over 140 characters")
		WriteErrorResponse(w, 400, "Chirp is too long")
		return
	}

	cleanedChirp := cleanChirp(params.Body)

	respBody := struct {
		Cleaned_body string `json:"cleaned_body"`
	}{
		Cleaned_body: cleanedChirp,
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		WriteErrorResponse(w, 500, "Something went wrong")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}
