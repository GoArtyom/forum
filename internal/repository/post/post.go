package post

import (
	"database/sql"

	"forum/internal/models"
)

type PostSqlite struct {
	db *sql.DB
}

func NewPostSqlite(db *sql.DB) *PostSqlite {
	return &PostSqlite{db: db}
}

func (r *PostSqlite) CreatePost(post *models.CreatePost) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	query := "INSERT INTO posts (title, content, user_id, user_name, create_at) VALUES ($1, $2, $3, $4, $5)"

	result, err := tx.Exec(query, post.Title, post.Content, post.UserId, post.UserName, post.CreateAt)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	postId, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	query2 := "INSERT INTO posts_categories (post_id, category_name) VALUES ($1, $2) "
	for _, category := range *post.Categories {
		_, err := tx.Exec(query2, postId, category)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	return int(postId), tx.Commit()
}
