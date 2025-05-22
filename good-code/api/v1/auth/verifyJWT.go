package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	middleware "github.com/chopstickleg/good-code/api/v1/_middleware"
	"github.com/dgrijalva/jwt-go"
)

func VerifyJWTHandler(w http.ResponseWriter, r *http.Request) {
	middleware.AllowMethods("GET")(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				http.Error(w, "Not authorized", http.StatusUnauthorized)
				return
			default:
				http.Error(w, "Internal server error loading cookie", http.StatusInternalServerError)
				return
			}
		}
		secretKey := os.Getenv("JWT_SECRET_KEY")
		if secretKey == "" {
			http.Error(w, "JWT_SECRET_KEY environment variable not set", http.StatusInternalServerError)
			return
		}
		err = verifyToken(cookie.Value, secretKey)
		if err != nil {
			switch {
			case errors.Is(err, fmt.Errorf("invalid token")):
				http.Error(w, "Not authorized", http.StatusUnauthorized)
				return
			default:
				http.Error(w, "Internal server error when verifying token"+err.Error(), http.StatusInternalServerError)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		response := json.NewEncoder(w)
		err = response.Encode(map[string]bool{"loggedIn": true})
		if err != nil {
			http.Error(w, "Failed to send response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})(w, r)
}

func verifyToken(tokenString string, secretKey string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
