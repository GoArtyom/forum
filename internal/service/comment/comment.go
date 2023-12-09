package comment

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type CommentService struct {
	repo repository.Comment
}

func NewCommentServer(repo repository.Comment) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) CreateComment(comment *models.CreateComment) error {
	return s.repo.CreateComment(comment)
}
