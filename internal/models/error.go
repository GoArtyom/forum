package models

import (
	"errors"
)

type DataForErr struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

const (
	UniqueEmail  = "UNIQUE constraint failed: users.email"
	UniqueName   = "UNIQUE constraint failed: users.name"
	IncorRequest = "FOREIGN KEY constraint failed"
	ErrEmail     = "Examples of valid email addresses: user@example.com, john.doe123@domain.co"
)

var (
	ErrUniqueUser = errors.New("unique user")
	ErrIncorData  = errors.New("incorrect password or email")
)
