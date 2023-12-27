package comment

import (
	"forum/internal/models"
	repo "forum/internal/repository"
)

type CommentService struct {
	repo repo.Comment
}

func NewCommentService(repo repo.Comment) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) CreateComment(comment *models.CreateComment) error {
	return s.repo.CreateComment(comment)
}

func (s *CommentService) GetAllCommentByPostId(postId int) ([]*models.Comment, error) {
	comments, err := s.repo.GetAllCommentByPostId(postId)
	if err != nil {
		return nil, err
	}
	for i, j := 0, len(comments)-1; i < j; i, j = i+1, j-1 {
		comments[i], comments[j] = comments[j], comments[i]
	}
	return comments, nil
}
