package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/adamararcane/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaUpgrade(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		WriteErrorResponse(w, 500, "Something went wrong")
		return
	}

	apikey, err := auth.GetApiKey(r.Header)
	if err != nil {
		log.Printf("User header is malformed: %s", err)
		WriteErrorResponse(w, 401, "ApiKey <token> apikey not found")
		return
	}

	if apikey != cfg.POLKA_KEY {
		log.Printf("ApiKey does not match")
		WriteErrorResponse(w, 401, "Invalid API Key")
		return
	}

	if params.Event != "user.upgraded" {
		WriteErrorResponse(w, 204, "Endpoint user.upgraded only")
		return
	}

	UserID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		log.Printf("Error parsing string UUID to UUID: %s", err)
		WriteErrorResponse(w, 400, "Invalid UUID")
		return
	}

	err = cfg.db.UpgradeToRedByID(r.Context(), UserID)
	if err != nil {
		log.Printf("User not found %s", err)
		WriteErrorResponse(w, 404, "User not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(204)
}
