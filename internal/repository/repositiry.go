package repository

import (
	"database/sql"

	"forum/internal/models"
	"forum/internal/repository/category"
	"forum/internal/repository/comment"
	commentvote "forum/internal/repository/commentVote"
	"forum/internal/repository/post"
	postvote "forum/internal/repository/postVote"
	"forum/internal/repository/session"
	"forum/internal/repository/user"
)

type User interface {
	CreateUser(user *models.CreateUser) error
	GetUserByEmail(email string) (*models.User, error)
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
	CreateSession(session *models.Session) error
	GetSessionByUserId(userId int) (*models.Session, error)
	GetSessionByUUID(uuid string) (*models.Session, error)
	DeleteSessionByUUID(uuid string) error
}

type Category interface {
	GetAllCategory() ([]*models.Category, error)
}

type PostVote interface {
	CreatePostVote(newVote *models.PostVote) error
	GetVoteByUserId(newVote *models.PostVote) (int, error)
	DeleteVoteByUserId(newVote *models.PostVote) error
}

type CommentVote interface {
	CreateCommentVote(newVote *models.CommentVote) error
}

type Repository struct {
	User
	Post
	Comment
	Session
	Category
	PostVote
	CommentVote
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User:        user.NewUserSqlite(db),
		Post:        post.NewPostSqlite(db),
		Comment:     comment.NewCommentSqlite(db),
		Session:     session.NewSessionSqlite(db),
		Category:    category.NewCategorySqlite(db),
		PostVote:    postvote.NewPostVoteSqlite(db),
		CommentVote: commentvote.NewCommentVoteSqlite(db),
	}
}
