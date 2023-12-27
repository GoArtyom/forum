package models

type UserInfoGoogle struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Sub   string `json:"sub"`
}

type UserInfoGitHub struct {
	Login  string `json:"login"`
	NodeId string `json:"node_id"`
}
