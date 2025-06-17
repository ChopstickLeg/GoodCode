package utils

import (
	"log"
	"net/http"

	db "github.com/chopstickleg/good-code/api/v1/_db"
	"github.com/google/go-github/v72/github"
)

func HandleInstallationTargetEvent(w http.ResponseWriter, body github.InstallationTargetEvent) {
	repository := body.GetRepository()
	if repository == nil {
		log.Printf("Installation target event missing repository data")
		http.Error(w, "Invalid installation target event: missing repository", http.StatusBadRequest)
		return
	}

	if repository.Owner == nil || repository.Owner.Login == nil {
		log.Printf("Installation target event missing repository owner data")
		http.Error(w, "Invalid installation target event: missing repository owner", http.StatusBadRequest)
		return
	}

	log.Printf("Processing installation target event for repository: %s", repository.GetFullName())

	conn, err := db.GetDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}

	err = conn.Model(&db.Repository{}).
		Where("id = ?", repository.GetID()).
		Update("owner", repository.Owner.GetLogin()).Error

	if err != nil {
		log.Printf("Failed to update repository owner for repo ID %d: %v", repository.GetID(), err)
		http.Error(w, "Failed to update repository owner", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully updated repository owner for repo ID %d", repository.GetID())
}
