package category

import (
	"forum/internal/models"
	repo "forum/internal/repository"
)

type CategoryService struct {
	repo *repo.Repository
}

func NewCategoryService(repo *repo.Repository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll() ([]*models.Category, error) {
	return s.repo.Category.GetAll()
}

func (s *CategoryService) Create(name string) error {
	return s.repo.Category.Create(name)
}

func (s *CategoryService) DeleteByName(name string) error {
	return s.repo.Category.DeleteByName(name)
}
