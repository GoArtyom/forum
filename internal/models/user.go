package models

const (
	LocalMode              = uint8(0)
	GoogleMode             = uint8(1)
	GitHubMode             = uint8(2)
	UserRole               = uint8(1)
	ModeratorRole          = uint8(8)
	ConsiderationModerator = uint8(5)
	AdminRole              = uint8(10)
)

type User struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Mode        uint8  `json:"mode"`
	Role        uint8  `json:"rols"`
	CountNotice int    `json:"count_notice"`
}

type CreateUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Mode     uint8  `json:"mode"`
}

type SignInUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Mode     uint8  `json:"mode"`
	Rols     uint8  `json:"rols"`
}

type UpdateRole struct {
	UserId   int   `json:"user_id"`
	UserRole uint8 `json:"user_role"`
	NewRole  uint8 `json:"new_role"`
}
