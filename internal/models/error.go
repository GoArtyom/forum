package models

import "errors"

const (
	UniqueEmail  = "UNIQUE constraint failed: users.email"
	UniqueName   = "UNIQUE constraint failed: users.name"
	IncorRequest = "FOREIGN KEY constraint failed"
)

var (
	UniqueUser = errors.New("unique user")
	IncorData  = errors.New("incorrect password or email")
)
