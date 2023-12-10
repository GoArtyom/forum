package post

import (
	"database/sql"
	"fmt"

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
		fmt.Println("Err1")
		return 0, err
	}

	query := "INSERT INTO posts (title, content, user_id, user_name, create_at) VALUES ($1, $2, $3, $4, $5)"

	result, err := tx.Exec(query, post.Title, post.Content, post.UserId, post.UserName, post.CreateAt)
	if err != nil {
		fmt.Println("Err2")
		tx.Rollback()
		return 0, err
	}

	postId, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Err3")
		tx.Rollback()
		return 0, err
	}

	query2 := "INSERT INTO posts_categories (post_id, category_name) VALUES ($1, $2) "
	for _, category := range *post.Categories {
		fmt.Printf("postId:%d category:%s\n", postId, category)
		_, err := tx.Exec(query2, postId, category)
		if err != nil {
			fmt.Println("Err4")
			tx.Rollback()
			return 0, err
		}
	}
	return int(postId), tx.Commit()
}

func (r *PostSqlite) GetPostById(postId int) (*models.Post, error) {
	post := &models.Post{}
	query := "SELECT * FROM posts WHERE id = $1"
	err := r.db.QueryRow(query, postId).Scan(&post.PostId, &post.Title, &post.Content, &post.UserId, &post.UserName, &post.CreateAt)
	return post, err
}
func (r *PostSqlite) GetAllPost() ([]*models.Post, error) {
	query := "SELECT * FROM posts"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	posts := make([]*models.Post, 0)
	for rows.Next() {
		post := new(models.Post)
		err := rows.Scan(&post.PostId, &post.Title, &post.Content,
			&post.UserId, &post.UserName, &post.CreateAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostSqlite) GetPostsByUserId(userId int) ([]*models.Post, error) {
	query := "SELECT * FROM posts WHERE user_id = $1"
	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	posts := make([]*models.Post, 0)
	for rows.Next() {
		post := new(models.Post)
		err := rows.Scan(&post.PostId, &post.Title, &post.Content,
			&post.UserId, &post.UserName, &post.CreateAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostSqlite) GetPostsByCategory(category string) ([]*models.Post, error) {
	query := "SELECT * FROM posts p JOIN posts_categories c ON p.id = c.post_id  WHERE c.category_name = $1"
	rows, err := r.db.Query(query, category)
	if err != nil {
		return nil, err
	}
	posts := make([]*models.Post, 0)
	for rows.Next() {
		post := new(models.Post)
		err := rows.Scan(&post.PostId, &post.Title, &post.Content,
			&post.UserId, &post.UserName, &post.CreateAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
