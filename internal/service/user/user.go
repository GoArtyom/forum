package user

import (
	"database/sql"
	"strings"

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
	// for Google
	if user.Mode == models.GoogleMode {
		space := strings.IndexRune(user.Name, ' ')
		user.Name = string(user.Name[0]) + "." + user.Name[space+1:]
	}
	// for Local
	if user.Mode == models.Local {
		user.Email = strings.ToLower(user.Email)
	}
	passwordHash := pkg.GetPasswordHash(user.Password)
	user.Password = passwordHash

	return s.repo.CreateUser(user)
}

func (s *UserService) SignInUser(user *models.SignInUser) (int, error) {
	// for Local
	if user.Mode == models.Local {
		user.Email = strings.ToLower(user.Email)
	}
	repoUser, err := s.repo.GetUserByEmail(user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, models.ErrIncorData
		} else {
			return 0, err
		}
	}
	if repoUser.Mode != models.Local {
		if user.Mode == models.GoogleMode {
			space := strings.IndexRune(user.Name, ' ')
			user.Name = string(user.Name[0]) + "." + user.Name[space+1:]
		}
		if user.Name != repoUser.Name {
			err = s.repo.UpdateUserNameById(repoUser.Id, user.Name)
			if err != nil {
				return 0, err
			}
		}
	}
	// for Local
	if user.Mode == models.Local {
		if repoUser.Password != pkg.GetPasswordHash(user.Password) {
			return 0, models.ErrIncorData
		}
	}
	return repoUser.Id, nil
}

func (s *UserService) GetUserByUserId(userId int) (*models.User, error) {
	return s.repo.GetUserByUserId(userId)
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s *UserService) UpdateUserNameById(userId int, newName string) error {
	return s.repo.UpdateUserNameById(userId, newName)
}
