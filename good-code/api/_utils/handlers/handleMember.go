package handlers

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/chopstickleg/good-code/api/_db"
	"gorm.io/gorm"

	"github.com/google/go-github/v72/github"
)

func HandleMemberEvent(w http.ResponseWriter, body github.MemberEvent) {
	action := body.GetAction()
	repository := body.GetRepo()
	member := body.GetMember()
	changes := body.GetChanges()

	log.Printf("Member event: %s for member: %s in repo: %s", action, member.GetLogin(), repository.GetFullName())

	conn, err := db.GetDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}

	switch action {
	case "added":
		if err := handleMemberAdded(conn, repository, member, changes); err != nil {
			log.Printf("Error handling member addition: %v", err)
			http.Error(w, "Failed to process member addition", http.StatusInternalServerError)
			return
		}
	case "edited":
		if err := handleMemberEdited(conn, repository, member, changes); err != nil {
			log.Printf("Error handling member edit: %v", err)
			http.Error(w, "Failed to process member edit", http.StatusInternalServerError)
			return
		}
	case "removed":
		if err := handleMemberRemoved(conn, repository, member); err != nil {
			log.Printf("Error handling member removal: %v", err)
			http.Error(w, "Failed to process member removal", http.StatusInternalServerError)
			return
		}
	default:
		log.Printf("Unhandled member action: %s", action)
	}
}

func handleMemberAdded(conn *gorm.DB, repository *github.Repository, member *github.User, changes *github.MemberChanges) error {
	count := int64(0)
	err := conn.Model(&db.UserLogin{}).
		Where(&db.UserLogin{GithubID: member.GetID()}).
		Count(&count).
		Error
	if err != nil {
		log.Printf("Failed to check user login for GitHub ID %d: %v", member.GetID(), err)
		return err
	}

	collaborator := db.UserRepositoryCollaborator{
		RepositoryID: repository.GetID(),
		GithubUserID: member.GetID(),
		GithubLogin:  member.GetLogin(),
		Role:         changes.Permission.GetTo(),
	}
	if count > 0 {
		var userLogin db.UserLogin
		err = conn.Where(&db.UserLogin{GithubID: member.GetID()}).
			First(&userLogin).
			Error
		if err != nil {
			return fmt.Errorf("failed to get user login ID: %v", err)
		}
		collaborator.UserLoginID = &userLogin.ID
	}
	err = conn.Create(&collaborator).Error
	return err
}

func handleMemberEdited(conn *gorm.DB, repository *github.Repository, member *github.User, changes *github.MemberChanges) error {
	return conn.Model(&db.UserRepositoryCollaborator{}).
		Where(&db.UserRepositoryCollaborator{GithubUserID: member.GetID(), RepositoryID: repository.GetID()}).
		Updates(&db.UserRepositoryCollaborator{Role: changes.Permission.GetTo()}).
		Error
}

func handleMemberRemoved(conn *gorm.DB, repository *github.Repository, member *github.User) error {
	return conn.Where(&db.UserRepositoryCollaborator{GithubUserID: member.GetID(), RepositoryID: repository.GetID()}).
		Delete(&db.UserRepositoryCollaborator{}).
		Error
}
