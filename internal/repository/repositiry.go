package repository

import (
	"database/sql"
	"forum/internal/models"
	"forum/internal/repository/user"
)

type User interface {
	CreateUser(user *models.User) error
}

type Post interface{}

type Comment interface{}

type Session interface{}

type Repository struct {
	User
	Post
	Comment
	Session
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User: user.NewUserSqlite3(db),
	}
}
