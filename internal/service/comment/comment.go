package comment

import (
	"forum/internal/models"
	repo "forum/internal/repository"
)

type CommentService struct {
	repo repo.Comment
}

func NewCommentServer(repo repo.Comment) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) CreateComment(comment *models.CreateComment) error {
	return s.repo.CreateComment(comment)
}

func (s *CommentService) GetAllCommentByPostId(postId int) ([]*models.Comment, error) {
	return s.repo.GetAllCommentByPostId(postId)
}
