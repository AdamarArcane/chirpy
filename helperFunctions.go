package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func WriteErrorResponse(w http.ResponseWriter, status int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	respBody := struct {
		Error string `json:"error"`
	}{
		Error: errorMsg,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling error response JSON: %s", err)
		w.Write([]byte("Something went really wrong"))
		return
	}
	w.Write(dat)
}

func cleanChirp(body string) string {
	bannedWords := []string{"kerfuffle", "sharbert", "fornax"}
	chirpWordList := strings.Split(body, " ")

	for i, word := range chirpWordList {
		for _, badWord := range bannedWords {
			if strings.EqualFold(badWord, word) {
				chirpWordList[i] = "****"
			}
		}
	}

	return strings.Join(chirpWordList, " ")
}
