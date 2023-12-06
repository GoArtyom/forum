package user

import (
	"database/sql"
	"forum/internal/models"
)

type UserSqlite struct {
	db *sql.DB
}

func NewUserSqlite3(db *sql.DB) *UserSqlite {
	return &UserSqlite{db: db}
}

func (r *UserSqlite) CreateUser(user *models.User) error {
	query := "INSERT INTO users (name, email, password_hash) VALUES($1, $2, $3)"
	_, err := r.db.Exec(query, user.Name, user.Email, user.Password)
	return err
}
