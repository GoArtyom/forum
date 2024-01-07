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

func (r *UserSqlite) Create(user *models.CreateUser) error {
	query := "INSERT INTO users (name, email, password_hash, mode, rols) VALUES($1, $2, $3, $4, $5)"
	_, err := r.db.Exec(query, user.Name, user.Email, user.Password, user.Mode, models.UserRole)

	return err
}

func (r *UserSqlite) GetByEmail(email string) (*models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE email = $1"
	err := r.db.QueryRow(query, email).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Mode, &user.Role)

	return &user, err
}

func (r *UserSqlite) GetById(userId int) (*models.User, error) {
	user := &models.User{}
	query := "SELECT * FROM users WHERE id = $1"
	err := r.db.QueryRow(query, userId).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Mode, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			err = models.ErrUniqueUser
		}
		return nil, err
	}

	return user, err
}

func (r *UserSqlite) GetAll() ([]*models.User, error) {
	query := "SELECT * FROM users"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	users := make([]*models.User, 0)
	for rows.Next() {
		user := new(models.User)
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Mode, &user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserSqlite) GetAllByRole(role uint8) ([]*models.User, error) {
	query := "SELECT * FROM users WHERE rols = ?"
	rows, err := r.db.Query(query, role)
	if err != nil {
		return nil, err
	}
	users := make([]*models.User, 0)
	for rows.Next() {
		user := new(models.User)
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Mode, &user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserSqlite) UpdateNameById(userId int, newName string) error {
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

func (r *UserSqlite) UpdateRoleById(userId int, newRole uint8) error {
	query := "UPDATE users SET rols = $1 WHERE id = $2"
	_, err := r.db.Exec(query, newRole, userId)
	return err
}
