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
	if err := conn.Model(&db.Repository{}).
		Where(&db.Repository{InstallationID: installation.GetID()}).
		Updates(&db.Repository{Enabled: false}).
		Error; err != nil {
		log.Printf("failed to fetch repository IDs for installation %d: %v", installation.GetID(), err)
		return err
	}
	return nil
}

func handleAppSuspended(conn *gorm.DB, installation *github.Installation) error {
	return conn.Model(&db.Repository{}).
		Where(&db.Repository{InstallationID: installation.GetID()}).
		Updates(db.Repository{Enabled: false}).
		Error
}

func handleAppUnsuspended(conn *gorm.DB, installation *github.Installation) error {
	return conn.Model(&db.Repository{}).
		Where(db.Repository{InstallationID: installation.GetID()}).
		Updates(&db.Repository{Enabled: true}).
		Error
}
func handleAppCreated(conn *gorm.DB, installation *github.Installation, repos []*github.Repository) error {
	for _, repo := range repos {
		var count int64
		err := conn.Model(&db.Repository{}).
			Where(&db.Repository{ID: repo.GetID()}).
			Count(&count).
			Error
		if err != nil {
			log.Printf("Failed to check repository %s: %v", repo.GetFullName(), err)
			return err
		}
		login, ownerID, collaborators, err := getRepoInfo(repo.GetID(), installation.GetID())
		if err != nil {
			log.Printf("Failed to get owner info for repo %s: %v", repo.GetFullName(), err)
			return err
		}
		if count == 0 {
			newRepo := db.Repository{
				ID:             repo.GetID(),
				Name:           repo.GetName(),
				Owner:          login,
				OwnerID:        ownerID,
				InstallationID: installation.GetID(),
			}
			if err := conn.Create(&newRepo).Error; err != nil {
				log.Printf("Failed to create repository record for %s: %v", repo.GetFullName(), err)
				return err
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
					RepositoryID:   newRepo.ID,
					GithubUserID:   collaborator.GetID(),
					GithubLogin:    collaborator.GetLogin(),
					Role:           collaborator.GetRoleName(),
					UserLoginID:    userLoginID,
					IsGoodCodeUser: userLoginID != nil,
				}
				if err := conn.Create(&collab).Error; err != nil {
					log.Printf("Failed to create collaborator record for %s in repo %s: %v", collaborator.GetLogin(), repo.GetFullName(), err)
					return err
				}
			}
		}
	}
	return nil
}

func getRepoInfo(repoId int64, installationID int64) (string, int64, []*github.User, error) {
	log.Printf("Using installation ID: %d", installationID)

	installationToken, err := utils.GetGitHubInstallationToken(installationID)
	if err != nil {
		log.Printf("Failed to get GitHub installation token: %v", err)
		return "", 0, nil, err
	}

	GHclient := github.NewClient(nil)
	authedGHClient := GHclient.WithAuthToken(installationToken)
	repo, _, err := authedGHClient.Repositories.GetByID(context.Background(), repoId)
	if err != nil {
		log.Printf("Failed to get repo info for repo ID %d: %v", repoId, err)
		return "", 0, nil, err
	}

	collaborators, _, err := authedGHClient.Repositories.ListCollaborators(context.Background(), repo.GetOwner().GetLogin(), repo.GetName(), nil)
	if err != nil {
		log.Printf("Failed to list collaborators for repo %s: %v", repo.GetFullName(), err)
		return "", 0, nil, err
	}
	return repo.GetOwner().GetLogin(), repo.GetOwner().GetID(), collaborators, nil
}
