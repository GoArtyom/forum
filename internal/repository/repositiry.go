package repository

import (
	"database/sql"

	"forum/internal/models"
	"forum/internal/repository/session"
	"forum/internal/repository/user"
)

type User interface {
	CreateUser(user *models.CreateUser) error
	GetUserByEmail(email string) (*models.User, error)
}

type Post interface{}

type Comment interface{}

type Session interface {
	CreateSession(session *models.Session) error
	GetSessionByUserId(userId int) (*models.Session, error)
	DeleteSessionByUUID(sessionId string) error
}

type Repository struct {
	User
	Post
	Comment
	Session
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User:    user.NewUserSqlite(db),
		Session: session.NewSessionSqlite(db),
	}
}
