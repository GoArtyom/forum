package user

import "forum/internal/repository"

type UserServise struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserServise {
	return &UserServise{repo: repo}
}
