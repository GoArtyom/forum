package service

import (
	"forum/internal/models"
	"forum/internal/repository"
	"forum/internal/service/session"
	"forum/internal/service/user"
)

type User interface {
	CreateUser(user *models.CreateUser) error
	SignInUser(user *models.SignInUser) (int, error)
}

type Post interface{}

type Comment interface{}

type Session interface {
	CreateSession(userId int) (*models.Session, error)
}

type Service struct {
	User
	Post
	Comment
	Session
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User:    user.NewUserService(repo),
		Session: session.NewSessionServise(repo),
	}
}
