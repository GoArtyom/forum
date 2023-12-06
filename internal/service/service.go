package service

import (
	"forum/internal/repository"
	"forum/internal/service/user"
)

type User interface{
	CreateUser()
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
