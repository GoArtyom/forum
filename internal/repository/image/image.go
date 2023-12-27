package image

import (
	"database/sql"
	"fmt"
	"forum/internal/models"
)

type ImageSqlite struct {
	db *sql.DB
}

func NewImageSqlite(db *sql.DB) *ImageSqlite {
	return &ImageSqlite{db: db}
}

func (r *ImageSqlite) CreateImageByPostId(newImage *models.CreateImage) error {
	query := "INSERT INTO posts_images (post_id, name, type) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(query, newImage.PostId, newImage.Name, newImage.Type)
	if err != nil {
		fmt.Println("Here")
		return err
	}
	return nil
}

func (r *ImageSqlite) DeleteImageByPostId(postId int) error {
	query := "DELETE FROM posts_images WHERE post_id = ?"
	_, err := r.db.Exec(query, postId)
	return err
}

func (r *ImageSqlite) GetImageByPostId(postId int) (*models.Image, error) {
	image := &models.Image{}

	query := "SELECT name, type FROM posts_images WHERE post_id = ?"
	err := r.db.QueryRow(query, postId).Scan(&image.Name, &image.Type)
	if err != nil {
		return nil, err
	}
	return image, nil
}
