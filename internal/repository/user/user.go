package user

import (
	"database/sql"

	"forum/internal/models"
)

type UserSqlite struct {
	db *sql.DB
}

func NewUserSqlite(db *sql.DB) *UserSqlite {
	return &UserSqlite{db: db}
}

func (r *UserSqlite) CreateUser(user *models.CreateUser) error {
	query := "INSERT INTO users (name, email, password_hash, mode) VALUES($1, $2, $3, $4)"
	_, err := r.db.Exec(query, user.Name, user.Email, user.Password, &user.Mode)

	return err
}

func (r *UserSqlite) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE email = $1"
	err := r.db.QueryRow(query, email).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Mode)

	return &user, err
}

func (r *UserSqlite) GetUserByUserId(userId int) (*models.User, error) {
	user := &models.User{}
	query := "SELECT * FROM users WHERE id = $1"
	err := r.db.QueryRow(query, userId).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Mode)
	return user, err
}

func (r *UserSqlite) UpdateUserNameById(userId int, newName string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	query := "UPDATE users SET name = $1 WHERE id = $2"
	_, err = tx.Exec(query, newName, userId)
	if err != nil {
		tx.Rollback()
		return err
	}
	query = "UPDATE posts SET user_name = $1 WHERE user_id = $2"
	_, err = tx.Exec(query, newName, userId)
	if err != nil {
		tx.Rollback()
		return err
	}
	query = "UPDATE comments SET user_name = $1 WHERE user_id = $2"
	_, err = tx.Exec(query, newName, userId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
