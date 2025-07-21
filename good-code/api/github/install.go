package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	db "github.com/chopstickleg/good-code/api/_db"
	middleware "github.com/chopstickleg/good-code/api/_middleware"
	utils "github.com/chopstickleg/good-code/api/_utils"
	"github.com/google/go-github/v72/github"
)

func HandleInstallationEvent(w http.ResponseWriter, r *http.Request) {
	middleware.AllowMethods("GET")(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				http.Error(w, "Not authorized", http.StatusUnauthorized)
				log.Printf("Error loading cookie: %v", err)
				return
			default:
				http.Error(w, "Internal server error loading cookie", http.StatusInternalServerError)
				log.Printf("Error loading cookie: %v", err)
				return
			}
		}
		userid, err := middleware.GetUserIDFromJWT(cookie.Value)
		if err != nil {
			log.Printf("Error verifying JWT: %v", err)
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		var body db.InstallationEvent
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			log.Printf("Error decoding request body: %v", err)
			return
		}

		if body.SetupAction != "created" {
			http.Error(w, "Unsupported action in installation event", http.StatusBadRequest)
			log.Printf("Unsupported action in installation event: %s", body.SetupAction)
			return
		}

		token, err := utils.GetGitHubInstallationToken(body.InstallationID)

		ghClient := github.NewClient(nil).WithAuthToken(token)

		installation, _, err := ghClient.Apps.GetInstallation(context.Background(), body.InstallationID)
		if err != nil {
			http.Error(w, "Failed to get installation details", http.StatusInternalServerError)
			log.Printf("Error getting installation details: %v", err)
			return
		}

		conn, err := db.GetDB()
		if err != nil {
			http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
			log.Printf("Error connecting to database: %v", err)
			return
		}

		err = conn.Model(&db.UserLogin{}).Updates(db.UserLogin{
			GithubID: installation.GetAccount().GetID(),
		}).
			Where(&db.UserLogin{ID: userid}).
			Error
		if err != nil {
			http.Error(w, "Failed to update user login", http.StatusInternalServerError)
			log.Printf("Error updating user login: %v", err)
			return
		}
	})(w, r)
}
