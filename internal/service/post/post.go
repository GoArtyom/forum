package post

import (
	"fmt"

	"forum/internal/models"
	repo "forum/internal/repository"
)

type PostService struct {
	repo repo.Post
	cat  repo.Category
}

func NewPostService(repo *repo.Repository) *PostService {
	return &PostService{
		repo: repo.Post,
		cat:  repo.Category,
	}
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
	err := s.cat.GetCategoryByName(category)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return s.repo.GetPostsByCategory(category)
}

func (s *PostService) GetPostsByLike(userId int) ([]*models.Post, error) {
	return s.repo.GetPostsByLike(userId)
}
