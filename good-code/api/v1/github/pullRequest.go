package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	db "github.com/chopstickleg/good-code/api/v1/_db"
	middleware "github.com/chopstickleg/good-code/api/v1/_middleware"
)

func PullRequestHandler(w http.ResponseWriter, r *http.Request) {
	middleware.AllowMethods(http.MethodGet)(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("auth")
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

		userId, err := middleware.GetUserIDFromJWT(token.Value)
		if err != nil {
			log.Printf("Error verifying JWT: %v", err)
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		conn, err := db.GetDB()
		if err != nil {
			http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
			return
		}

		var user db.UserLogin
		err = conn.Where("id = ?", userId).First(&user).Error
		if err != nil {
			log.Printf("Error finding user %d: %v", userId, err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		var aiRoasts []db.AiRoast
		err = conn.Joins("Repository").
			Where("Repository.owner_id = ? AND Repository.enabled = ?", user.GithubId, true).
			Find(&aiRoasts).Error

		if err != nil {
			log.Printf("Error retrieving AI roasts for user %d: %v", userId, err)
			http.Error(w, "Error retrieving data from database", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(aiRoasts)
		if err != nil {
			http.Error(w, "Error sending response", http.StatusInternalServerError)
			return
		}
	})(w, r)
}
