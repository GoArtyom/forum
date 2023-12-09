package category

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type CategoryService struct {
	repo repository.Category
}

func NewCategoryService(repo repository.Category) *CategoryService {
	return &CategoryService{repo: repo}
}
func (s *CategoryService) GetAllCategory() ([]*models.Category, error) {
	return s.repo.GetAllCategory()
}
