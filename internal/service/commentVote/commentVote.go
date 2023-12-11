package commentvote

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type CommentVoteService struct {
	repo repository.CommentVote
}

func NewCommentVoteService(repo repository.CommentVote) *CommentVoteService {
	return &CommentVoteService{repo: repo}
}

func (s *CommentVoteService) CreateCommentVote(newVote *models.CommentVote) error {
	return nil
}
