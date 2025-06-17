package db

type UserLogin struct {
	ID       int
	Email    string
	Password []byte
	Name     string
	Enabled  bool
	GithubId int64
}
type AiRoast struct {
	ID            int
	AiAnalysis    string
	RepoId        int64
	PullRequestId int64
}
type PullRequest struct {
	ID int64
	Repository
	Number int
}
type Repository struct {
	ID             int64
	Name           string
	Owner          string
	OwnerId        int64
	Enabled        bool
	InstallationId int64
}
