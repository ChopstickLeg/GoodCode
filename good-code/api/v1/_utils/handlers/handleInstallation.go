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
	case "created", "new_permissions_accepted":
		if err := handleAppCreated(conn, installation, body.Repositories); err != nil {
			log.Printf("Error handling app creation: %v", err)
			http.Error(w, "Failed to process app creation", http.StatusInternalServerError)
			return
		}
	default:
		log.Printf("Unhandled installation action: %s", action)
	}
}

func handleAppUninstalled(conn *gorm.DB, installation *github.Installation) error {
	var repoIDs []int64
	if err := conn.Model(&db.Repository{}).
		Where("installation_id = ?", installation.GetID()).
		Pluck("id", &repoIDs).
		Error; err != nil {
		log.Printf("failed to fetch repository IDs for installation %d: %v", installation.GetID(), err)
		return err
	}

	for _, repoID := range repoIDs {
		if err := conn.Where("repo_id = ?", repoID).Delete(&db.AiRoast{}).Error; err != nil {
			log.Printf("failed to delete AI roasts for repository %d: %w", repoID, err)
			return err
		}

		if err := conn.Where("id = ?", repoID).Delete(&db.Repository{}).Error; err != nil {
			log.Printf("failed to delete repository %d: %w", repoID, err)
			return err
		}
	}
	return nil
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
func handleAppCreated(conn *gorm.DB, installation *github.Installation, repos []*github.Repository) error {
	for _, repo := range repos {
		var count int64
		err := conn.Model(&db.Repository{}).
			Where("repo_id = ?", repo.GetID()).
			Count(&count).
			Error
		if err != nil {
			log.Printf("Failed to check repository %s: %v", repo.GetFullName(), err)
			return err
		}
		if count == 0 {
			newRepo := db.Repository{
				ID:             repo.GetID(),
				Name:           repo.GetName(),
				Owner:          repo.Owner.GetLogin(),
				OwnerID:        repo.Owner.GetID(),
				InstallationID: installation.GetID(),
			}
			if err := conn.Create(&newRepo).Error; err != nil {
				log.Printf("Failed to create repository record for %s: %v", repo.GetFullName(), err)
				return err
			}
		}
	}
	return nil
}
