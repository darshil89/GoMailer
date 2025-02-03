package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"mailingService/models"

	"github.com/golang-jwt/jwt/v5"
)

// helper function to respond with error
func RespondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// helper function to respond with success
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func GetToken(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid format")
		return
	}
	if creds.Secret != os.Getenv("SECRET") {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorised")
		return
	}

	expirationTime := time.Hour
	claims := &jwt.RegisteredClaims{
		Issuer:    "mailservice",
		Subject:   "user",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	var jwtKey = []byte(os.Getenv("JWT_KEY"))

	log.Printf("JWT key in getToken : %v", jwtKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}
	response := struct {
		Token string `json:"token"`
	}{
		Token: signedToken,
	}
	RespondWithJSON(w, http.StatusAccepted, response)
}
