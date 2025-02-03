package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"mailingService/models"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))

// helper function to respond with error
func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// helper function to respond with success
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Mailing Service is up and running")
	})
	r.Post("/api/getToken", GetToken)
	r.Post("/api/sendEmail", authenticate(validateEmail(SendEmail)))

	fmt.Println("Server listening on port :8080")
	log.Fatal(http.ListenAndServeTLS(":8080", "server.crt", "server.key", r))
}

func authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// signedToken := r.Header.Get("Authorization")
		var data models.EmailData
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		log.Printf("Data in authenticate is : %v", data)

		// token, err := jwt.ParseWithClaims(signedToken, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		// 	return jwtKey, nil
		// })
		// if err != nil || !token.Valid {
		// 	respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		// 	return
		// }

		ctx := context.WithValue(r.Context(), "emailData", data)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}

func validateEmail(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data models.EmailData
		data, ok := r.Context().Value("emailData").(models.EmailData)

		log.Printf("Data in validation is : %v", data)
		if !ok {
			respondWithError(w, http.StatusInternalServerError, "Failed to retrieve email data in validateEmail")
			return
		}

		verifier := emailverifier.NewVerifier().DisableDomainSuggest().DisableCatchAllCheck()
		ret, err := verifier.Verify(data.Email)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Email address not verified")
			log.Printf("verify email address failed, error is: %v", err)
			return
		}

		if !ret.Syntax.Valid {
			respondWithError(w, http.StatusBadRequest, "Invalid email syntax")
			return
		}
		if !ret.HasMxRecords {
			respondWithError(w, http.StatusBadRequest, "Email domain does not have MX records")
			return
		}
		if ret.Reachable == "no" {
			respondWithError(w, http.StatusBadRequest, "Email address is not reachable")
			return
		}
		if ret.RoleAccount {
			respondWithError(w, http.StatusBadRequest, "Role-based email addresses are not allowed")
			return
		}

		ctx := context.WithValue(r.Context(), "emailData", data)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}
