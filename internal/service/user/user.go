package user

import (
	"errors"

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
	user.Password = passwordHash
	return s.repo.CreateUser(user)
}

func (s *UserServise) SignInUser(user *models.SignInUser) (int, error) {
	repoUser, err := s.repo.GetUserByEmail(user.Email)
	if err != nil {
		return 0, err
	}

	if repoUser.Password != pkg.GetPasswordHash(user.Password) {
		return 0, errors.New("incorrect password")
	}
	return repoUser.Id, nil
}
