package db

type User_login struct {
	ID       int
	Email    string
	Password []byte
	Name     string
	Enabled  bool
}
type Pull_request struct {
	ID            int
	Author_id     int
	Author_name   string
	Source_branch string
	Target_branch string
	Has_comments  bool
	AIComments    string
}

type AI_request struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	System string `json:"system"`
	Stream bool   `json:"stream"`
}