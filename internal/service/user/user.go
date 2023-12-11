package user

import (
	"forum/internal/models"
	repo "forum/internal/repository"
	"forum/pkg"
)

type UserService struct {
	repo repo.User
}

func NewUserService(repo repo.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user *models.CreateUser) error {
	passwordHash := pkg.GetPasswordHash(user.Password)
	user.Password = passwordHash
	return s.repo.CreateUser(user)
}

func (s *UserService) SignInUser(user *models.SignInUser) (int, error) {
	repoUser, err := s.repo.GetUserByEmail(user.Email)
	if err != nil {
		return 0, models.IncorData
	}

	if repoUser.Password != pkg.GetPasswordHash(user.Password) {
		return 0, models.IncorData
	}
	return repoUser.Id, nil
}

func (s *UserService) GetUserByUserId(userId int) (*models.User, error) {
	return s.repo.GetUserByUserId(userId)
}
