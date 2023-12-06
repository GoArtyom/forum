package repository

import "database/sql"

type User interface{}

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
	return &Repository{}
}
