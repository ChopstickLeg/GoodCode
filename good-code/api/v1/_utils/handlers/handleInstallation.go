package handlers

import (
	"log"
	"net/http"

	db "github.com/chopstickleg/good-code/api/v1/_db"
	"github.com/google/go-github/v72/github"
	"gorm.io/gorm"
)

func HandleInstallationEvent(w http.ResponseWriter, body github.InstallationEvent) {
	action := body.GetAction()
	installation := body.GetInstallation()

	log.Printf("Installation event: %s for installation ID: %d", action, installation.GetID())

	conn, err := db.GetDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}

	switch action {
	case "deleted":
		if err := handleAppUninstalled(conn, installation); err != nil {
			log.Printf("Error handling app uninstallation: %v", err)
			http.Error(w, "Failed to process app uninstallation", http.StatusInternalServerError)
			return
		}
	case "suspend":
		if err := handleAppSuspended(conn, installation); err != nil {
			log.Printf("Error handling app suspension: %v", err)
			http.Error(w, "Failed to process app suspension", http.StatusInternalServerError)
			return
		}
	case "unsuspend":
		if err := handleAppUnsuspended(conn, installation); err != nil {
			log.Printf("Error handling app unsuspension: %v", err)
			http.Error(w, "Failed to process app unsuspension", http.StatusInternalServerError)
			return
		}
	default:
		log.Printf("Unhandled installation action: %s", action)
	}
}

func handleAppUninstalled(conn *gorm.DB, installation *github.Installation) error {
	return conn.Model(&db.Repository{}).
		Where("installation_id = ?", installation.GetID()).
		Update("enabled", false).Error
}

func handleAppSuspended(conn *gorm.DB, installation *github.Installation) error {
	return conn.Model(&db.Repository{}).
		Where("installation_id = ?", installation.GetID()).
		Update("enabled", false).Error
}

func handleAppUnsuspended(conn *gorm.DB, installation *github.Installation) error {
	return conn.Model(&db.Repository{}).
		Where("installation_id = ?", installation.GetID()).
		Update("enabled", true).Error
}
