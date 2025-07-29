package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	db "github.com/chopstickleg/good-code/api/_db"
	middleware "github.com/chopstickleg/good-code/api/_middleware"
	repository "github.com/chopstickleg/good-code/api/_utils/repository"
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

		repoId, err := repository.GetRepoId(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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

		var roasts []db.AiRoast
		err = conn.
			Omit("Repository").
			Where(&db.AiRoast{RepoID: repoId}).
			Find(&roasts).
			Error
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
