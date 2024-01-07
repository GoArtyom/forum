package service

import (
	"forum/internal/models"
	"forum/internal/repository"
	"forum/internal/service/category"
	"forum/internal/service/comment"
	commentvote "forum/internal/service/commentVote"
	"forum/internal/service/image"
	"forum/internal/service/notification"
	"forum/internal/service/post"
	postvote "forum/internal/service/postVote"
	"forum/internal/service/report"
	"forum/internal/service/session"
	"forum/internal/service/user"
)

type User interface {
	Create(user *models.CreateUser) error
	SignIn(user *models.SignInUser) (int, error)
	GetById(userId int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetAll() ([]*models.User, error)
	GetAllByRole(role uint8) ([]*models.User, error)
	UpdateNameById(userId int, newName string) error
	UpdateRoleById(upRole *models.UpdateRole) error
	FilterByRole(allUsers []*models.User, role uint8) []*models.User
}

type Post interface {
	Create(post *models.CreatePost) (int, error)
	GetById(postId int) (*models.Post, error)
	GetAll() ([]*models.Post, error)
	GetAllByUserVote(userId int) ([]*models.Post, error)
	GetAllByUserId(userId int) ([]*models.Post, error)
	GetByCategory(category string) ([]*models.Post, error)
	GetByVote(userId, vote int) ([]*models.Post, error)
	UpdateById(upPost *models.UpdatePost) error
	DeleteById(formDel *models.DeletePost) error
}

type Comment interface {
	Create(comment *models.CreateComment) (int, error)
	GetById(commentId int) (*models.Comment, error)
	GetAllByPostId(postId int) ([]*models.Comment, error)
	GetAllByUserId(userId int) ([]*models.Comment, error)
	GetAllByUserVote(userId int) ([]*models.Comment, error)
	UpdateById(upComment *models.UpdateComment) error
	DeleteById(formDel *models.DeleteComment) error
}

type Session interface {
	Create(userId int) (*models.Session, error)
	GetByUUID(uuid string) (*models.Session, error)
	DeleteByUUID(uuid string) error
}
type Category interface {
	Create(name string) error
	GetAll() ([]*models.Category, error)
	DeleteByName(name string) error
}

type PostVote interface {
	Create(newVote *models.PostVote) (uint8, error)
}

type CommentVote interface {
	Create(newVote *models.CommentVote) (uint8, error)
}
type Image interface {
	CreateByPostId(newImage *models.CreateImage) error
	GetByPostId(postId int) (*models.Image, error)
	DeleteByPostId(PostId int) error
}

type Report interface {
	Create(report *models.CreateReport) error
	DeleteById(reportId, resp int) error
	GetAll() ([]*models.Report, error)
}

type Notification interface {
	Create(notic *models.Notification) error
	GetAllByAuthorId(authorId int) ([]*models.Notification, error)
	GetCountByAuthorId(authorId int) (int, error)
	Delete(formDel *models.DeleteNotification) error
}

type Service struct {
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

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User:         user.NewUserService(repo.User),
		Post:         post.NewPostService(repo),
		Comment:      comment.NewCommentService(repo.Comment),
		Session:      session.NewSessionService(repo.Session),
		Category:     category.NewCategoryService(repo),
		PostVote:     postvote.NewPostVoteService(repo.PostVote),
		CommentVote:  commentvote.NewCommentVoteService(repo.CommentVote),
		Image:        image.NewImageService(repo.Image),
		Report:       report.NewReportService(repo),
		Notification: notification.NewNoticService(repo),
	}
}
