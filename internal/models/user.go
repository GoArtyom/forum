package models

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
