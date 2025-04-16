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
	Comments      Comment
}
type Comment struct {
	ID          int
	Author_id   int
	Author_name string
	Text        string
}
