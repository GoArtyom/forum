package models

import "errors"

const (
	UniqueEmail  = "UNIQUE constraint failed: users.email"
	UniqueName   = "UNIQUE constraint failed: users.name"
	IncorRequest = "FOREIGN KEY constraint failed"
	ErrPassword  = "Please enter a password that meets the following requirements:\n" +
		"Contains at least one lowercase letter\n" +
		"Contains at least one uppercase letter\n" +
		"Contains at least one number\n" +
		"Contains at least one of the following special characters: !@#$%^&*"
	ErrEmail = "Incorrect email.\n" +
		"Examples of valid email addresses: user@example.com, john.doe123@domain.co"
)

var (
	UniqueUser = errors.New("unique user")
	IncorData  = errors.New("incorrect password or email")
)
