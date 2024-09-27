package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
	}

	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	timeNowUTC := time.Now().UTC()
	expDur := (time.Second * 3600)
	expTime := timeNowUTC.Add(expDur)
	subject := userID.String()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(timeNowUTC),
		ExpiresAt: jwt.NewNumericDate(expTime),
		Subject:   subject,
	})

	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	tknStruct, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		log.Printf("Error validating token: %s", err)
		return uuid.UUID{}, err
	}

	strUUID, err := tknStruct.Claims.GetSubject()
	if err != nil {
		log.Printf("Error getting subject from token claims: %s", err)
		return uuid.UUID{}, err
	}

	UUID, err := uuid.Parse(strUUID)
	if err != nil {
		log.Printf("Error parsing UUID string to UUID: %s", err)
		return uuid.UUID{}, err
	}

	return UUID, nil
}

func MakeRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("error generating refresh token: %w", err)
	}

	rfToken := hex.EncodeToString(bytes)
	return rfToken, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("header does not have correct prefix")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	return token, nil
}
