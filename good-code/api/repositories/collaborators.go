package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	db "github.com/chopstickleg/good-code/api/_db"
	middleware "github.com/chopstickleg/good-code/api/_middleware"
	utils "github.com/chopstickleg/good-code/api/_utils"
	repository "github.com/chopstickleg/good-code/api/_utils/repository"
	"github.com/google/go-github/v72/github"
)

func GetCollaboratorsHandler(w http.ResponseWriter, r *http.Request) {
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

		repoIdStr := r.URL.Query().Get("repoId")
		repoId, err := strconv.ParseInt(repoIdStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid repository ID", http.StatusBadRequest)
			return
		}

		conn, err := db.GetDB()
		if err != nil {
			http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
			return
		}

		var user db.UserLogin
		err = conn.Where(&db.UserLogin{ID: userId}).First(&user).Error
		if err != nil {
			log.Printf("Error retrieving user with ID %d: %v", userId, err)
			http.Error(w, "Error retrieving user from database", http.StatusInternalServerError)
			return
		}

		hasAccess, err := repository.GetRepoAccess(repoId, user, conn)
		if err != nil {
			log.Printf("Error checking repository access for user %d and repo %d: %v", userId, repoId, err)
			http.Error(w, "Error checking repository access", http.StatusInternalServerError)
			return
		}

		if !hasAccess {
			http.Error(w, "Not authorized to access this repository", http.StatusForbidden)
			return
		}

		var collaborators []db.UserRepositoryCollaborator
		err = conn.
			Omit("Repository").
			Where(&db.UserRepositoryCollaborator{RepositoryID: repoId}).
			Find(&collaborators).
			Error
		if err != nil {
			log.Printf("Error retrieving collaborators for repo %d: %v", repoId, err)
			http.Error(w, "Error retrieving data from database", http.StatusInternalServerError)
			return
		}

		var installationId int64
		err = conn.Model(&db.UserLogin{}).
			Select("user_logins.installation_id").
			Joins("JOIN repositories ON repositories.owner_id = user_logins.id").
			Where("repositories.id = ?", repoId).
			Scan(&installationId).
			Error

		if err != nil {
			log.Printf("Error retrieving installation ID for repo %d: %v", repoId, err)
			http.Error(w, "Error retrieving installation ID", http.StatusInternalServerError)
			return
		}

		installationToken, err := utils.GetGitHubInstallationToken(installationId)
		if err != nil {
			log.Printf("Failed to get GitHub installation token: %v", err)
			http.Error(w, "Unable to get GitHub installation token", http.StatusInternalServerError)
			return
		}

		GHclient := github.NewClient(nil)
		authedGHClient := GHclient.WithAuthToken(installationToken)

		for i := range collaborators {
			collaborator := &collaborators[i]
			if collaborator.UserLoginID != nil {
				user, _, err := authedGHClient.Users.GetByID(context.Background(), collaborator.GithubUserID)
				if err != nil {
					log.Printf("Failed to get user by ID: %v", err)
					continue
				}
				collaborator.GithubAvatarURL = user.GetAvatarURL()
			}
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(collaborators)
		if err != nil {
			http.Error(w, "Error sending response", http.StatusInternalServerError)
			return
		}
	})(w, r)
}
