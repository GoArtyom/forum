package service

import "forum/internal/repository"

type User interface{}

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
	return &Service{}
}
