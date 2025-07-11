package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	db "github.com/chopstickleg/good-code/api/_db"
	middleware "github.com/chopstickleg/good-code/api/_middleware"
)

func GetRoastsHandler(w http.ResponseWriter, r *http.Request) {
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

		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		var repoIdStr string
		for i, part := range pathParts {
			if part == "repositories" && i+1 < len(pathParts) {
				repoIdStr = pathParts[i+1]
				break
			}
		}

		if repoIdStr == "" {
			http.Error(w, "Repository ID not found in URL", http.StatusBadRequest)
			return
		}

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

		var repo db.Repository
		err = conn.Where(&db.Repository{ID: repoId}).First(&repo).Error
		if err != nil {
			http.Error(w, "Repository not found", http.StatusNotFound)
			return
		}

		var owner db.UserLogin
		err = conn.Where(&db.UserLogin{ID: userId}).First(&owner).Error
		if err != nil {
			log.Printf("Error retrieving owner with ID %d: %v", userId, err)
			http.Error(w, "Error retrieving owner from database", http.StatusInternalServerError)
			return
		}

		isOwner := repo.OwnerID == owner.GithubID

		var isCollaborator bool
		if !isOwner {
			var collaborator db.UserRepositoryCollaborator
			err = conn.Where(&db.UserRepositoryCollaborator{RepositoryID: repoId, UserLoginID: &userId}).First(&collaborator).Error
			if err == nil {
				isCollaborator = true
			}
		}

		if !isOwner && !isCollaborator {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		var roasts []db.AiRoast
		err = conn.Where(&db.AiRoast{RepoID: repoId}).Find(&roasts).Error
		if err != nil {
			log.Printf("Error retrieving AI roasts for repo %d: %v", repoId, err)
			http.Error(w, "Error retrieving data from database", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(roasts)
		if err != nil {
			http.Error(w, "Error sending response", http.StatusInternalServerError)
			return
		}
	})(w, r)
}
