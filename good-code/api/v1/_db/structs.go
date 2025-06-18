package db

type UserLogin struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"-"`
	Name     string `json:"name"`
	GithubId int64  `json:"github_id"`
	Enabled  bool   `json:"enabled"`

	// Repositories they own
	OwnedRepositories []Repository `gorm:"foreignKey:OwnerId;references:GithubId" json:"owned_repositories"`

	// Repositories they collaborate on (many-to-many)
	CollaboratingRepositories []Repository `gorm:"many2many:user_repository_collaborators;" json:"collaborating_repositories"`
}

type Repository struct {
	ID             int64  `gorm:"primaryKey" json:"id"`
	Name           string `json:"name"`
	Owner          string `json:"owner"`
	OwnerId        int64  `json:"owner_id"`
	Enabled        bool   `json:"enabled"`
	InstallationID int64  `json:"installation_id"`

	// Owner relationship
	OwnerUser UserLogin `gorm:"foreignKey:OwnerId;references:GithubId" json:"owner_user"`

	// Collaborators (many-to-many)
	Collaborators []UserLogin `gorm:"many2many:user_repository_collaborators;" json:"collaborators"`

	// AI Roasts
	AiRoasts []AiRoast `gorm:"foreignKey:RepoId" json:"ai_roasts"`
}

type UserRepositoryCollaborator struct {
	UserLoginID  int64  `gorm:"primaryKey" json:"user_login_id"`
	RepositoryID int64  `gorm:"primaryKey" json:"repository_id"`
	Role         string `json:"role"`

	UserLogin  UserLogin  `gorm:"foreignKey:UserLoginID" json:"user_login"`
	Repository Repository `gorm:"foreignKey:RepositoryID" json:"repository"`
}

type AiRoast struct {
	ID                int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	RepoId            int64  `json:"repo_id"`
	PullRequestNumber int    `json:"pull_request_number"`
	Content           string `json:"content"`

	Repository Repository `gorm:"foreignKey:RepoId"`
}
