package user

import (
	"database/sql"
	"strings"

	"forum/internal/models"
	repo "forum/internal/repository"
	"forum/pkg"
)

type UserService struct {
	user repo.User
}

func NewUserService(repo repo.User) *UserService {
	return &UserService{user: repo}
}

func (s *UserService) Create(user *models.CreateUser) error {
	// for Google
	if user.Mode == models.GoogleMode {
		space := strings.IndexRune(user.Name, ' ')
		user.Name = string(user.Name[0]) + "." + user.Name[space+1:]
	}
	// for Local
	if user.Mode == models.LocalMode {
		user.Email = strings.ToLower(user.Email)
	}
	passwordHash := pkg.GetPasswordHash(user.Password)
	user.Password = passwordHash

	return s.user.Create(user)
}

func (s *UserService) SignIn(user *models.SignInUser) (int, error) {
	// for Local
	if user.Mode == models.LocalMode {
		user.Email = strings.ToLower(user.Email)
	}
	repoUser, err := s.user.GetByEmail(user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, models.ErrIncorData
		} else {
			return 0, err
		}
	}
	if repoUser.Mode != models.LocalMode {
		if user.Mode == models.GoogleMode {
			space := strings.IndexRune(user.Name, ' ')
			user.Name = string(user.Name[0]) + "." + user.Name[space+1:]
		}
		if user.Name != repoUser.Name {
			err = s.user.UpdateNameById(repoUser.Id, user.Name)
			if err != nil {
				return 0, err
			}
		}
	}
	// for Local
	if user.Mode == models.LocalMode {
		if repoUser.Password != pkg.GetPasswordHash(user.Password) {
			return 0, models.ErrIncorData
		}
	}
	return repoUser.Id, nil
}

func (s *UserService) GetById(userId int) (*models.User, error) {
	return s.user.GetById(userId)
}

func (s *UserService) GetByEmail(email string) (*models.User, error) {
	return s.user.GetByEmail(email)
}

func (s *UserService) UpdateNameById(userId int, newName string) error {
	return s.user.UpdateNameById(userId, newName)
}

func (s *UserService) UpdateRoleById(upRole *models.UpdateRole) error {

	switch upRole.UserRole {
	case models.UserRole:
		if upRole.NewRole != models.ModeratorRole && upRole.NewRole != models.ConsiderationModerator {
			return models.ErrUpdateRole
		}
	case models.ConsiderationModerator:
		if upRole.NewRole != models.ModeratorRole && upRole.NewRole != models.UserRole {
			return models.ErrUpdateRole
		}
	case models.ModeratorRole:
		if upRole.NewRole != models.UserRole {
			return models.ErrUpdateRole
		}
	default:
		return models.ErrUpdateRole
	}

	err := s.user.UpdateRoleById(upRole.UserId, upRole.NewRole)
	if err != nil && err == sql.ErrNoRows {
		err = models.ErrUpdateRole
	}

	return err
}

func (s *UserService) GetAll() ([]*models.User, error) {
	users, err := s.user.GetAll()
	if err != nil {
		return nil, err
	}
	return users[1:], nil
}

func (s *UserService) GetAllByRole(role uint8) ([]*models.User, error) {
	return s.user.GetAllByRole(role)
}

func (s *UserService) FilterByRole(allUsers []*models.User, role uint8) []*models.User {
	var res []*models.User
	for _, user := range allUsers {
		if user.Role == role {
			res = append(res, user)
		}
	}
	return res
}
