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

func (r *CategorySqlite) GetAllCategory() ([]*models.Category, error) {
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
