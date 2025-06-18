package handlers

import (
	"log"
	"net/http"

	db "github.com/chopstickleg/good-code/api/v1/_db"
	"gorm.io/gorm"

	"github.com/google/go-github/v72/github"
)

func HandleMemberEvent(w http.ResponseWriter, body github.MemberEvent) {
	action := body.GetAction()
	repository := body.GetRepo()
	member := body.GetMember()

	log.Printf("Member event: %s for member: %s in repo: %s", action, member.GetLogin(), repository.GetFullName())

	conn, err := db.GetDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}

	switch action {
	case "added":
		if err := handleMemberAdded(conn, repository, member); err != nil {
			log.Printf("Error handling member addition: %v", err)
			http.Error(w, "Failed to process member addition", http.StatusInternalServerError)
			return
		}
	default:
		log.Printf("Unhandled member action: %s", action)
	}
}

func handleMemberAdded(conn *gorm.DB, repository *github.Repository, member *github.User) error {
	
}
