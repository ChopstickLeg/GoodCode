package repository

import (
	"fmt"

	db "github.com/chopstickleg/good-code/api/_db"

	"gorm.io/gorm"
)

func GetRepoAccess(repoId int64, owner db.UserLogin, conn *gorm.DB) (bool, error) {
	var repo db.Repository
	err := conn.Where(&db.Repository{ID: repoId}).First(&repo).Error
	if err != nil {
		return false, fmt.Errorf("error retrieving repository with ID %d: %w", repoId, err)
	}

	isOwner := repo.OwnerID == owner.GithubID

	var isCollaborator bool
	if !isOwner {
		var collaborator db.UserRepositoryCollaborator
		err = conn.Where(&db.UserRepositoryCollaborator{RepositoryID: repoId, UserLoginID: &owner.ID}).First(&collaborator).Error
		if err == nil {
			isCollaborator = true
		}
	}
	return (isOwner || isCollaborator), nil
}
