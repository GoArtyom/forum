package post

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type PostService struct {
	repo repository.Post
}

func NewPostService(repo repository.Post) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePost(post *models.CreatePost) (int, error) {
	return s.repo.CreatePost(post)
}

func (s *PostService) GetPostById(postId int) (*models.Post, error) {
	return s.repo.GetPostById(postId)
}

func (s *PostService) GetAllPost() ([]*models.Post, error) {
	return s.repo.GetAllPost()
}

func (s *PostService) GetPostsByUserId(userId int) ([]*models.Post, error) {
	return s.repo.GetPostsByUserId(userId)
}

func (s *PostService) GetPostsByCategory(category string) ([]*models.Post, error) {
	return s.repo.GetPostsByCategory(category)
}

func (s *PostService) GetPostsByLike(userId int) ([]*models.Post, error) {
	return s.repo.GetPostsByLike(userId)
}
