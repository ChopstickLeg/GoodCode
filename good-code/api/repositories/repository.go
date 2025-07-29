package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	db "github.com/chopstickleg/good-code/api/_db"
	middleware "github.com/chopstickleg/good-code/api/_middleware"
	repository "github.com/chopstickleg/good-code/api/_utils/repository"
)

func GetRepoHandler(w http.ResponseWriter, r *http.Request) {
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

		var user db.UserLogin
		userId, err := middleware.GetUserIDFromJWT(token.Value)
		if err != nil {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		repoId, err := repository.GetRepoId(r)
		if err != nil {
			http.Error(w, "Invalid repository ID", http.StatusBadRequest)
			return
		}

		conn, err := db.GetDB()
		if err != nil {
			http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
			return
		}

		err = conn.Where(&db.UserLogin{ID: userId}).First(&user).Error
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		hasAccess, err := repository.GetRepoAccess(repoId, user, conn)
		if err != nil {
			http.Error(w, "Error checking repository access", http.StatusInternalServerError)
			return
		}
		if !hasAccess {
			http.Error(w, "Not authorized to access this repository", http.StatusUnauthorized)
			return
		}

		var repo db.Repository
		err = conn.
			Preload("OwnerUser").
			Preload("Collaborators").
			Preload("AiRoasts").
			Where(&db.Repository{ID: repoId}).
			First(&repo).
			Error
		if err != nil {
			http.Error(w, "Repository not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(repo); err != nil {
			http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
			return
		}

	})(w, r)
}
