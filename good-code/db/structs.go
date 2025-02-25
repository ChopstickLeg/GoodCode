package db

type User_login struct {
	ID       int
	Email    string
	Password []byte
	Enabled  bool
}
