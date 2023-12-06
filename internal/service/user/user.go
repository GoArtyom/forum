package user

import (
	"forum/internal/models"
	"forum/internal/repository"
	"forum/pkg"
)

type UserServise struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserServise {
	return &UserServise{repo: repo}
}

func (s *UserServise) CreateUser(user *models.CreateUser) error {
	passwordHash := pkg.GetPasswordHash(user.Password)
	userS := &models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: passwordHash,
	}
	return s.repo.CreateUser(userS)
}
