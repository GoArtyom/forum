package service

import (
	"forum/internal/models"
	"forum/internal/repository"
	"forum/internal/service/post"
	"forum/internal/service/session"
	"forum/internal/service/user"
)

type User interface {
	CreateUser(user *models.CreateUser) error
	SignInUser(user *models.SignInUser) (int, error)
	GetUserByUserId(userId int) (*models.User, error)
}

type Post interface {
	CreatePost(post *models.CreatePost) (int, error)
}

type Comment interface{}

type Session interface {
	CreateSession(userId int) (*models.Session, error)
	GetSessionByUUID(uuid string) (*models.Session, error)
	DeleteSessionByUUID(uuid string) error
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
		Post:    post.NewPostService(repo),
		Session: session.NewSessionServise(repo),
	}
}
