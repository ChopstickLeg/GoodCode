package db

type UserLogin struct {
	ID       int
	Email    string
	Password []byte
	Name     string
	Enabled  bool
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
	ID      int64
	Name    string
	Owner   string
	Enabled bool
}