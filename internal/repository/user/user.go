package user

import (
	"database/sql"
	"fmt"

	"forum/internal/models"
)

type UserSqlite struct {
	db *sql.DB
}

func NewUserSqlite(db *sql.DB) *UserSqlite {
	return &UserSqlite{db: db}
}

func (r *UserSqlite) CreateUser(user *models.CreateUser) error {
	query := "INSERT INTO users (name, email, password_hash) VALUES($1, $2, $3)"
	_, err := r.db.Exec(query, user.Name, user.Email, user.Password)
	if err != nil {
		if err.Error() == models.UniqueEmail || err.Error() == models.UniqueName {
			return models.UniqueUser
		}
		return err
	}
	return err
}

func (r *UserSqlite) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	fmt.Println("repo:", email)
	query := "SELECT * FROM users WHERE email = $1"
	err := r.db.QueryRow(query, email).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	return &user, err
}

func (r *UserSqlite) GetUserByUserId(userId int) (*models.User, error) {
	user := &models.User{}
	fmt.Println("repo:", userId)
	query := "SELECT * FROM users WHERE id = $1"
	err := r.db.QueryRow(query, userId).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	return user, err
}
