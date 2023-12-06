package service

import (
	"forum/internal/models"
	"forum/internal/repository"
	"forum/internal/service/user"
)

type User interface {
	CreateUser(user *models.CreateUser) error
}

type Post interface{}

type Comment interface{}

type Session interface{}

type Service struct {
	User
	Post
	Comment
	Session
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User: user.NewUserService(repo),
	}
}
