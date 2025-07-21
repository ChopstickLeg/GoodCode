package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	middleware "github.com/chopstickleg/good-code/api/_middleware"
	authentication "github.com/chopstickleg/good-code/api/_utils/authentication"
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
		err = authentication.VerifyToken(cookie.Value)
		if err != nil {
			switch {
			case errors.Is(err, fmt.Errorf("invalid token")):
				http.Error(w, "Not authorized", http.StatusUnauthorized)
				return
			default:
				http.Error(w, "Internal server error when verifying token: "+err.Error(), http.StatusInternalServerError)
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
