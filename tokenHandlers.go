package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/adamararcane/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshJWT(w http.ResponseWriter, r *http.Request) {

	rftoken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting auth token: %s", err)
		WriteErrorResponse(w, 401, "Token expired or does not exist")
		return
	}

	revoked_at, err := cfg.db.RevokeCheck(r.Context(), rftoken)
	if err != nil {
		log.Printf("Error checking revoked status: %s", err)
		WriteErrorResponse(w, 500, "Error checking if token is revoked")
		return
	}

	if revoked_at.Valid == true {
		log.Printf("Token has been revoked")
		WriteErrorResponse(w, 401, "Token has been revoked")
		return
	}

	UUID, err := cfg.db.GetUserFromRefreshToken(r.Context(), rftoken)
	if err != nil {
		log.Printf("Error getting UUID from token: %s", err)
		WriteErrorResponse(w, 401, "Token does not exist")
		return
	}

	JWT, err := auth.MakeJWT(UUID, cfg.JWT_SECRET)
	if err != nil {
		log.Printf("Error making JWT: %s", err)
		WriteErrorResponse(w, 500, "Error making new JWT")
		return
	}

	resp := RefreshResp{
		Token: JWT,
	}

	dat, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		WriteErrorResponse(w, 500, "Something went wrong")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}

func (cfg *apiConfig) handlerRevokeRFToken(w http.ResponseWriter, r *http.Request) {
	rftoken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting auth token: %s", err)
		WriteErrorResponse(w, 401, "Token expired or does not exist")
		return
	}

	err = cfg.db.RevokeToken(r.Context(), rftoken)
	if err != nil {
		log.Printf("Error revoking RFToken: %s", err)
		WriteErrorResponse(w, 500, "Error revoking token")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(204)
}
