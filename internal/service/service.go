package service

import (
	"forum/internal/models"
	"forum/internal/repository"
	"forum/internal/service/category"
	"forum/internal/service/comment"
	commentvote "forum/internal/service/commentVote"
	"forum/internal/service/post"
	postvote "forum/internal/service/postVote"
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
	GetPostsByUserId(userId int) ([]*models.Post, error)
	GetPostsByCategory(category string) ([]*models.Post, error)
	GetPostsByLike(userId int) ([]*models.Post, error)
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

type PostVote interface {
	CreatePostVote(newVote *models.PostVote) error
}

type CommentVote interface {
	CreateCommentVote(newComment *models.CommentVote) error
}

type Service struct {
	User
	Post
	Comment
	Session
	Category
	PostVote
	CommentVote
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User:        user.NewUserService(repo),
		Post:        post.NewPostService(repo),
		Comment:     comment.NewCommentServer(repo),
		Session:     session.NewSessionService(repo),
		Category:    category.NewCategoryService(repo),
		PostVote:    postvote.NewPostVoteService(repo),
		CommentVote: commentvote.NewCommentVoteService(repo),
	}
}
