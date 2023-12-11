package commentvote

import (
	"database/sql"

	"forum/internal/models"
	repo "forum/internal/repository"
)

type CommentVoteService struct {
	repo repo.CommentVote
}

func NewCommentVoteService(repo repo.CommentVote) *CommentVoteService {
	return &CommentVoteService{repo: repo}
}

func (s *CommentVoteService) CreateCommentVote(newVote *models.CommentVote) error {
	vote, err := s.repo.GetVoteByUserId(newVote)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if vote != 0 { // проверяем наличие vote
		err = s.repo.DeleteVoteByUserId(newVote)
		if err != nil {
			return err
		}
	}
	if vote != newVote.Vote {
		err = s.repo.CreateCommentVote(newVote)
		if err != nil {
			return err
		}
	}
	return nil
}
