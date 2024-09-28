package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/adamararcane/chirpy/internal/auth"
	"github.com/adamararcane/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		WriteErrorResponse(w, 500, "Something went wrong")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		WriteErrorResponse(w, 500, "Error hashing password")
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{Email: params.Email, HashedPassword: hashedPassword})
	if err != nil {
		log.Printf("Error creating user: %s", err)
		WriteErrorResponse(w, 500, "User at email already exists")
		return
	}

	response := UserResp{
		ID:         user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email,
	}

	dat, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		WriteErrorResponse(w, 500, "Something went wrong")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(dat)
}

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		WriteErrorResponse(w, 500, "Something went wrong")
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error getting user (user dne): %s", err)
		WriteErrorResponse(w, 401, "Incorrect email or password")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		log.Printf("User password and hashedPW do not match: %s", err)
		WriteErrorResponse(w, 401, "Incorrect email or password")
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.JWT_SECRET)
	if err != nil {
		log.Printf("Error making JWT: %s", err)
		WriteErrorResponse(w, 500, "Error making JWT")
		return
	}

	rfToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Error making refresh token: %s", err)
		WriteErrorResponse(w, 500, "Error making refresh token")
		return
	}

	rfTokenItem, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{Token: rfToken, UserID: user.ID})
	if err != nil {
		log.Printf("Error creating refresh token DB Item: %s", err)
		WriteErrorResponse(w, 500, "Error adding refresh token to DB")
		return
	}

	response := UserWithToken{
		ID:            user.ID,
		Created_at:    user.CreatedAt,
		Updated_at:    user.UpdatedAt,
		Email:         user.Email,
		Token:         token,
		Refresh_token: rfTokenItem.Token,
	}

	dat, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		WriteErrorResponse(w, 500, "Something went wrong")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
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

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		WriteErrorResponse(w, 500, "Something went wrong")
		return
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("User header is malformed: %s", err)
		WriteErrorResponse(w, 401, "Bearer <token> authentication token not found")
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.JWT_SECRET)
	if err != nil {
		log.Printf("User token is invalid or expried")
		WriteErrorResponse(w, 401, "Token is malformed or missing")
		return
	}

	newHashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing new user password: %s", err)
		WriteErrorResponse(w, 500, "Error hashing new user password")
		return
	}

	updateUser := database.UpdateUserEmailAndPasswordParams{
		Email:          params.Email,
		HashedPassword: newHashedPassword,
		ID:             userID,
	}

	updatedUser, err := cfg.db.UpdateUserEmailAndPassword(r.Context(), updateUser)

	resp := UserResp{
		ID:         updatedUser.ID,
		Created_at: updatedUser.CreatedAt,
		Updated_at: updatedUser.UpdatedAt,
		Email:      updatedUser.Email,
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
