package handlers

import (
	"log"
	"net/http"

	db "github.com/chopstickleg/good-code/api/_db"
	"github.com/google/go-github/v72/github"
	"gorm.io/gorm"
)

func HandleRepositoryEvent(w http.ResponseWriter, body github.RepositoryEvent) {
	action := body.GetAction()
	repository := body.GetRepo()

	log.Printf("Repository event: %s for repo: %s", action, repository.GetFullName())

	conn, err := db.GetDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}

	switch action {
	case "deleted":
		if err := handleRepositoryDeleted(conn, repository); err != nil {
			log.Printf("Error handling repository deletion: %v", err)
			http.Error(w, "Failed to process repository deletion", http.StatusInternalServerError)
			return
		}
	case "transferred":
		if err := handleRepositoryTransferred(conn, body); err != nil {
			log.Printf("Error handling repository transfer: %v", err)
			http.Error(w, "Failed to process repository transfer", http.StatusInternalServerError)
			return
		}
	case "renamed":
		if err := handleRepositoryRenamed(conn, body); err != nil {
			log.Printf("Error handling repository rename: %v", err)
			http.Error(w, "Failed to process repository rename", http.StatusInternalServerError)
			return
		}
	case "created":
		if err := handleRepositoryCreated(conn, body); err != nil {
			log.Printf("Error handling repository creation: %v", err)
			http.Error(w, "Failed to process repository creation", http.StatusInternalServerError)
			return
		}
	default:
		log.Printf("Unhandled repository action: %s", action)
	}
}

func handleRepositoryDeleted(conn *gorm.DB, repository *github.Repository) error {
	repoID := repository.GetID()

	if err := conn.Where(&db.Repository{ID: repoID}).
		Delete(&db.AiRoast{}).Error; err != nil {
		log.Printf("failed to delete AI roasts for repository %d: %w", repoID, err)
		return err
	}

	if err := conn.Where(&db.Repository{ID: repoID}).
		Delete(&db.Repository{}).Error; err != nil {
		log.Printf("failed to delete repository %d: %w", repoID, err)
		return err
	}

	return nil
}

func handleRepositoryTransferred(conn *gorm.DB, body github.RepositoryEvent) error {
	repository := body.GetRepo()

	return conn.Model(&db.Repository{}).
		Where(&db.Repository{ID: repository.GetID()}).
		Updates(&db.Repository{Owner: repository.GetOwner().GetLogin(), OwnerID: repository.GetOwner().GetID()}).
		Error
}

func handleRepositoryRenamed(conn *gorm.DB, body github.RepositoryEvent) error {
	repository := body.GetRepo()
	changes := body.GetChanges()

	log.Printf("Repository renamed from %s to %s",
		changes.GetRepo().Name.GetFrom(), repository.GetName())

	return conn.Model(&db.Repository{}).
		Where(&db.Repository{ID: repository.GetID()}).
		Updates(&db.Repository{Name: repository.GetName()}).
		Error
}

func handleRepositoryCreated(conn *gorm.DB, body github.RepositoryEvent) error {
	repository := body.GetRepo()

	repo := db.Repository{
		ID:    repository.GetID(),
		Name:  repository.GetName(),
		Owner: repository.GetOwner().GetLogin(),
	}
	return conn.Create(&repo).Error
}
