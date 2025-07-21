package handlers

import (
	"context"
	"log"
	"net/http"

	db "github.com/chopstickleg/good-code/api/_db"
	utils "github.com/chopstickleg/good-code/api/_utils"
	"github.com/google/go-github/v72/github"
	"gorm.io/gorm"
)

func HandleRepositoryInstallationEvent(w http.ResponseWriter, body github.InstallationRepositoriesEvent) {
	action := body.GetAction()
	log.Printf("Installation event: %s for installation ID: %d", action, body.GetInstallation().GetID())
	switch action {
	case "added":
		handleRepositoryAdded(w, body)
	case "removed":
		handleRepositoryRemoved(w, body)
	default:
		log.Printf("Unhandled installation action: %s", action)
		http.Error(w, "Unhandled installation action", http.StatusBadRequest)
		return
	}
}

func handleRepositoryAdded(w http.ResponseWriter, body github.InstallationRepositoriesEvent) {
	installationId := body.Installation.GetID()

	installationToken, err := utils.GetGitHubInstallationToken(installationId)
	if err != nil {
		log.Printf("Failed to get GitHub installation token: %v", err)
		http.Error(w, "Unable to get GitHub installation token", http.StatusInternalServerError)
		return
	}

	GHclient := github.NewClient(nil)
	authedGHClient := GHclient.WithAuthToken(installationToken)

	conn, err := db.GetDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	for _, repo := range body.RepositoriesAdded {
		fullRepo, _, err := authedGHClient.Repositories.GetByID(context.Background(), repo.GetID())
		if err != nil {
			log.Printf("Failed to get repository by ID: %v", err)
			continue
		}
		err = conn.Create(&db.Repository{
			ID:             fullRepo.GetID(),
			Name:           fullRepo.GetName(),
			Owner:          fullRepo.GetOwner().GetLogin(),
			OwnerID:        fullRepo.GetOwner().GetID(),
			InstallationID: body.GetInstallation().GetID(),
		}).Error
		if err != nil {
			http.Error(w, "Failed to add repository to database", http.StatusInternalServerError)
			return
		}

		collaborators, _, err := authedGHClient.Repositories.ListCollaborators(context.Background(), fullRepo.GetOwner().GetLogin(), fullRepo.GetName(), nil)
		if err != nil {
			log.Printf("Failed to list collaborators for repo %s: %v", fullRepo.GetName(), err)
			continue
		}
		for _, collaborator := range collaborators {
			var userLoginID *int64
			var userLogin db.UserLogin
			err = conn.Where(&db.UserLogin{GithubID: collaborator.GetID()}).First(&userLogin).Error
			if gorm.ErrRecordNotFound == err {
				log.Printf("User login not found for collaborator %s in repo %s", collaborator.GetLogin(), repo.GetFullName())
				userLoginID = nil
			} else if err != nil {
				log.Printf("Failed to find user login for collaborator %s in repo %s: %v", collaborator.GetLogin(), repo.GetFullName(), err)
				continue
			} else {
				userLoginID = &userLogin.ID
			}
			log.Printf("Creating collaborator record for %s in repo %s", collaborator.GetLogin(), repo.GetFullName())
			collab := db.UserRepositoryCollaborator{
				RepositoryID:   fullRepo.GetID(),
				GithubUserID:   collaborator.GetID(),
				GithubLogin:    collaborator.GetLogin(),
				Role:           collaborator.GetRoleName(),
				UserLoginID:    userLoginID,
				IsGoodCodeUser: userLoginID != nil,
			}
			if err := conn.Create(&collab).Error; err != nil {
				log.Printf("Failed to create collaborator record for %s in repo %s: %v", collaborator.GetLogin(), repo.GetFullName(), err)
				http.Error(w, "Failed to create collaborator record", http.StatusInternalServerError)
				return
			}
		}
	}
}

func handleRepositoryRemoved(w http.ResponseWriter, body github.InstallationRepositoriesEvent) {
	conn, err := db.GetDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	for _, repo := range body.RepositoriesRemoved {
		err = conn.Where(&db.Repository{ID: repo.GetID()}).
			Updates(&db.Repository{Enabled: false}).Error
		if err != nil {
			http.Error(w, "Failed to remove repository from database", http.StatusInternalServerError)
			return
		}
	}
}
