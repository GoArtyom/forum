package category

import (
	"forum/internal/models"
	repo "forum/internal/repository"
)

type CategoryService struct {
	repo repo.Category
}

func NewCategoryService(repo repo.Category) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAllCategory() ([]*models.Category, error) {
	return s.repo.GetAllCategory()
}
