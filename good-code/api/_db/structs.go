package db

import "time"

type UserLogin struct {
	ID             int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Email          string `json:"email"`
	Password       []byte `json:"-"`
	Name           string `json:"name"`
	GithubID       int64  `gorm:"uniqueIndex" json:"github_id"`
	InstallationID int64  `json:"installation_id"`
	Enabled        bool   `json:"enabled"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Repositories they own
	OwnedRepositories []Repository `gorm:"foreignKey:OwnerID;references:GithubID" json:"owned_repositories"`

	// Repositories they collaborate on (many-to-many)
	CollaboratingRepositories []Repository `gorm:"many2many:user_repository_collaborators;" json:"collaborating_repositories"`
}

type Repository struct {
	ID      int64  `gorm:"primaryKey" json:"id"`
	Name    string `json:"name"`
	Owner   string `json:"owner"`
	OwnerID int64  `json:"owner_id"`
	Enabled bool   `gorm:"default:true" json:"enabled"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Owner relationship
	OwnerUser UserLogin `gorm:"foreignKey:OwnerID;references:GithubID" json:"owner_user"`

	// Collaborators
	Collaborators []UserRepositoryCollaborator `gorm:"foreignKey:RepositoryID" json:"collaborators"`

	// AI Roasts
	AiRoasts []AiRoast `gorm:"foreignKey:RepoID" json:"ai_roasts"`
}

type UserRepositoryCollaborator struct {
	ID           int64 `gorm:"primaryKey;autoIncrement" json:"id"`
	RepositoryID int64 `json:"repository_id"`

	GithubUserID    int64  `json:"github_user_id"`
	GithubLogin     string `json:"github_login"`
	Role            string `json:"role"`
	IsGoodCodeUser  bool   `json:"is_good_code_user"`
	GithubAvatarURL string `gorm:"-" json:"github_avatar_url"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	UserLoginID *int64     `gorm:"default:null" json:"user_login_id,omitempty"`
	UserLogin   *UserLogin `gorm:"foreignKey:UserLoginID" json:"user_login,omitempty"`
	Repository  Repository `gorm:"foreignKey:RepositoryID" json:"-"`
}

type AiRoast struct {
	ID                int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	RepoID            int64     `json:"repo_id"`
	PullRequestNumber int       `json:"pull_request_number"`
	Content           string    `json:"content"`
	IsOpen            bool      `json:"is_open"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Repository Repository `gorm:"foreignKey:RepoID" json:"-"`
}

type InstallationEvent struct {
	InstallationID int64  `json:"installation_id"`
	SetupAction    string `json:"setup_action"`
}
