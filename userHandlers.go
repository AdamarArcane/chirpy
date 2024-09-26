package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		WriteErrorResponse(w, 500, "Something went wrong")
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)

	response := User{
		ID:         user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email,
	}

	dat, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		WriteErrorResponse(w, 500, "Something went wrong")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(dat)
}

func (cfg *apiConfig) handlerResetUsers(w http.ResponseWriter, r *http.Request) {
	if cfg.PLATFORM != "dev" {
		log.Printf(".env PLATFORM not set to dev")
		WriteErrorResponse(w, 403, "Reset is forbidden")
		return
	} else {
		err := cfg.db.ResetUsers(r.Context())
		if err != nil {
			log.Printf("Error resetting users database: %s", err)
			WriteErrorResponse(w, 500, "Error resetting users database")
			return
		}
		fmt.Println("Users database has been reset")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("Users database has been reset"))
	}
}
