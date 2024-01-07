package category

import (
	"database/sql"

	"forum/internal/models"
)

type CategorySqlite struct {
	db *sql.DB
}

func NewCategorySqlite(db *sql.DB) *CategorySqlite {
	return &CategorySqlite{db: db}
}

func (r *CategorySqlite) GetAll() ([]*models.Category, error) {
	query := "SELECT * FROM category"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	categories := make([]*models.Category, 0)
	for rows.Next() {
		category := new(models.Category)
		err := rows.Scan(&category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategorySqlite) GetByName(name string) (string, error) {
	query := "SELECT name FROM category WHERE name = ?"
	var res string
	err := r.db.QueryRow(query, name).Scan(&res)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (r *CategorySqlite) Create(name string) error {
	query := "INSERT OR IGNORE INTO category (name) VALUES (?)"
	_, err := r.db.Exec(query, name)

	return err
}

func (r *CategorySqlite) DeleteByName(name string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	query := "DELETE FROM category WHERE name = ?"
	res, err := tx.Exec(query, name)
	res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}
	countDel, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if countDel == 0 {
		tx.Rollback()
		return models.ErrCategory
	}

	query = "DELETE FROM posts WHERE id NOT IN (SELECT post_id FROM posts_categories)"

	_, err = tx.Exec(query)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return models.ErrCategory
		}
		return err
	}

	return tx.Commit()
}
