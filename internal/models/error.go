package models

import (
	"errors"
)

type DataForErr struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Report   string `json:"report"`
}

const (
	UniqueEmail  = "UNIQUE constraint failed: users.email"
	UniqueName   = "UNIQUE constraint failed: users.name"
	IncorRequest = "FOREIGN KEY constraint failed"
)

// text err for field
const (
	TextErrEmail         = "Examples of valid email addresses: user@example.com, john.doe123@domain.co"
	TextErrFieldTooLong  = "This field is too long. Maximum is %d characters"
	TextErrFieldTooShort = "This field is too short. Minimum is %d characters"
	TextErrEmptyField    = "This field is empty"
	TextErrExtraSpace    = "This field has extra spaces"
)

// text errors for password
const (
	TextErrSpecialChars = "Contains at least one of the following special characters: !@#$%^&*"
	TextErrOneNumber    = "Contains at least one number"
	TextErrUpLetter     = "Contains at least one uppercase letter"
	TextErrLowLetter    = "Contains at least one lowercase letter"
)

// text errors for image
const (
	TextErrSize     = "File size exceeds 20 MB"
	TextErrTypeFile = "Invalid file type .%s (use JPEG, JPG, PNG, GIF)"
)

var (
	ErrUniqueUser   = errors.New("unique user")
	ErrIncorData    = errors.New("incorrect password or email")
	ErrUser         = errors.New("invalid request to user")
	ErrUpdateRole   = errors.New("invalid request to update user role")
	ErrPost         = errors.New("invalid request to post")
	ErrComment      = errors.New("invalid request to comment")
	ErrCategory     = errors.New("invalid request to category")
	ErrReport       = errors.New("invalid request to report")
	ErrNotification = errors.New("invalid request to notification")
)
