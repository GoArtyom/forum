package repository

import (
	"database/sql"

	"forum/internal/models"
	"forum/internal/repository/category"
	"forum/internal/repository/comment"
	commentvote "forum/internal/repository/commentVote"
	"forum/internal/repository/image"
	"forum/internal/repository/notification"
	"forum/internal/repository/post"
	postvote "forum/internal/repository/postVote"
	"forum/internal/repository/report"
	"forum/internal/repository/session"
	"forum/internal/repository/user"
)

type User interface {
	Create(user *models.CreateUser) error
	GetByEmail(email string) (*models.User, error)
	GetById(userId int) (*models.User, error)
	GetAll() ([]*models.User, error)
	GetAllByRole(role uint8) ([]*models.User, error)
	UpdateNameById(userId int, newName string) error
	UpdateRoleById(userId int, newRole uint8) error
}

type Post interface {
	Create(post *models.CreatePost) (int, error)
	GetById(postId int) (*models.Post, error)
	GetAll() ([]*models.Post, error)
	GetAllByUserId(userId int) ([]*models.Post, error)
	GetByCategory(category string) ([]*models.Post, error)
	GetByVote(userId, vote int) ([]*models.Post, error)
	UpdateById(upPost *models.UpdatePost) error
	DeleteById(postId int) error
	DeleteByCategory(category string) error
}

type Comment interface {
	Create(comment *models.CreateComment) (int, error)
	GetById(commentId int) (*models.Comment, error)
	GetByVote(userId, vote int) ([]*models.Comment, error)
	GetAllByPostId(postId int) ([]*models.Comment, error)
	GetAllByUserId(userId int) ([]*models.Comment, error)
	UpdateById(upComment *models.UpdateComment) error
	DeleteById(commentId int) error
}

type Session interface {
	Create(session *models.Session) error
	GetByUserId(userId int) (*models.Session, error)
	GetByUUID(uuid string) (*models.Session, error)
	DeleteByUUID(uuid string) error
}

type Category interface {
	Create(name string) error
	GetAll() ([]*models.Category, error)
	GetByName(name string) (string, error)
	DeleteByName(name string) error
}

type PostVote interface {
	Create(newVote *models.PostVote) error
	GetByUserId(newVote *models.PostVote) (int, error)
	DeleteByUserId(newVote *models.PostVote) error
}

type CommentVote interface {
	Create(newVote *models.CommentVote) error
	GetByUserId(newVote *models.CommentVote) (int, error)
	DeleteByUserId(newVote *models.CommentVote) error
}

type Image interface {
	CreateByPostId(newImage *models.CreateImage) error
	DeleteByPostId(postId int) error
	GetByPostId(postId int) (*models.Image, error)
}

type Report interface {
	Create(report *models.CreateReport) error
	GetById(reportId int) (*models.Report, error)
	GetAll() ([]*models.Report, error)
	DeleteById(reportId int) error
}

type Notification interface {
	Create(notic *models.Notification) error
	GetAllByAuthorId(authorId int) ([]*models.Notification, error)
	GetCountByAuthorId(authorId int) (int, error)
	DeleteById(noticId int) error
	DeleteByAuthorId(authorId int) error
	DeleteByUserId(formDel *models.DeleteNotification) error
}

type Repository struct {
	User
	Post
	Comment
	Session
	Category
	PostVote
	CommentVote
	Image
	Report
	Notification
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User:         user.NewUserSqlite(db),
		Post:         post.NewPostSqlite(db),
		Comment:      comment.NewCommentSqlite(db),
		Session:      session.NewSessionSqlite(db),
		Category:     category.NewCategorySqlite(db),
		PostVote:     postvote.NewPostVoteSqlite(db),
		CommentVote:  commentvote.NewCommentVoteSqlite(db),
		Image:        image.NewImageSqlite(db),
		Report:       report.NewReportSqlite(db),
		Notification: notification.NewNoticService(db),
	}
}
