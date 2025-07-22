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
	handlers "github.com/chopstickleg/good-code/api/_utils/handlers"
	"github.com/google/go-github/v72/github"
)

func HandleInstallationEvent(w http.ResponseWriter, r *http.Request) {
	middleware.AllowMethods("POST")(func(w http.ResponseWriter, r *http.Request) {
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

		if body.SetupAction != "install" {
			http.Error(w, "Unsupported action in installation event", http.StatusBadRequest)
			log.Printf("Unsupported action in installation event: %s", body.SetupAction)
			return
		}

		token, err := utils.GenerateGitHubJWT()

		ghClientJWT := github.NewClient(nil).WithAuthToken(token)

		installation, _, err := ghClientJWT.Apps.GetInstallation(context.Background(), body.InstallationID)
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

		err = conn.Model(&db.UserLogin{}).
			Where(&db.UserLogin{ID: userid}).
			Updates(db.UserLogin{
				GithubID: installation.GetID(),
			}).
			Error
		if err != nil {
			http.Error(w, "Failed to update user login", http.StatusInternalServerError)
			log.Printf("Error updating user login: %v", err)
			return
		}

		installationToken, err := utils.GetGitHubInstallationToken(installation.GetID())
		if err != nil {
			http.Error(w, "Failed to get GitHub installation token", http.StatusInternalServerError)
			log.Printf("Error getting GitHub installation token: %v", err)
			return
		}
		ghClient := github.NewClient(nil).WithAuthToken(installationToken)

		repositories, _, err := ghClient.Apps.ListRepos(context.Background(), &github.ListOptions{})

		err = handlers.HandleAppCreated(conn, installation, repositories.Repositories)
		if err != nil {
			http.Error(w, "Failed to handle app creation", http.StatusInternalServerError)
			log.Printf("Error handling app creation: %v", err)
			return
		}
	})(w, r)
}
