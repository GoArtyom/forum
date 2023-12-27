package models

const (
	Local      = uint8(0)
	GoogleMode = uint8(1)
	GitHubMode = uint8(2)
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Mode     uint8
}

type CreateUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Mode     uint8
}

type SignInUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Mode     uint8
}
