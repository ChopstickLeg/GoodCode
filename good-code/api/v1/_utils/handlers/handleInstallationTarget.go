package handlers

import (
	"log"
	"net/http"

	db "github.com/chopstickleg/good-code/api/v1/_db"
	"github.com/google/go-github/v72/github"
)

func HandleInstallationTargetEvent(w http.ResponseWriter, body github.InstallationTargetEvent) {
	installation := body.GetInstallation()
	repository := body.GetRepository()

	if installation == nil || repository == nil || repository.GetOwner() == nil {
		log.Printf("Installation target event missing installation, repository, or owner data")
		http.Error(w, "Invalid installation target event: missing installation, repository, or owner", http.StatusBadRequest)
		return
	}

	log.Printf("Processing installation target event for installation: %d", installation.GetID())

	conn, err := db.GetDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}

	err = conn.Model(&db.Repository{}).
		Where(&db.Repository{InstallationID: installation.GetID()}).
		Updates(&db.Repository{Owner: repository.GetOwner().GetLogin()}).
		Error

	if err != nil {
		log.Printf("Failed to update repository owner for installation ID %d: %v", installation.GetID(), err)
		http.Error(w, "Failed to update repository owner", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully updated repository owner for repo ID %d", repository.GetID())
}
