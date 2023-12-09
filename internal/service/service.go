package service

import (
	"forum/internal/models"
	"forum/internal/repository"
	"forum/internal/service/category"
	"forum/internal/service/comment"
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
	GetPostById(postId int) (*models.Post, error)
	GetAllPost() ([]*models.Post, error)
}

type Comment interface {
	CreateComment(comment *models.CreateComment) error
	GetAllCommentByPostId(postId int) ([]*models.Comment, error)
}

type Session interface {
	CreateSession(userId int) (*models.Session, error)
	GetSessionByUUID(uuid string) (*models.Session, error)
	DeleteSessionByUUID(uuid string) error
}
type Category interface {
	GetAllCategory() ([]*models.Category, error)
}

type Service struct {
	User
	Post
	Comment
	Session
	Category
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User:     user.NewUserService(repo),
		Post:     post.NewPostService(repo),
		Comment:  comment.NewCommentServer(repo),
		Session:  session.NewSessionService(repo),
		Category: category.NewCategoryService(repo),
	}
}
