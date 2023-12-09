package post

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type PostServise struct {
	repo repository.Post
}

func NewPostService(repo repository.Post) *PostServise {
	return &PostServise{repo: repo}
}

func (s *PostServise) CreatePost(post *models.CreatePost) (int, error) {
	return s.repo.CreatePost(post)
}

func (s *PostServise) GetPostById(postId int) (*models.Post, error) {
	return s.repo.GetPostById(postId)
}
