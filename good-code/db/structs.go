package db

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       int
	Email    string
	Password []byte
	Enabled  bool
}
