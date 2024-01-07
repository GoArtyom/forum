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

func (s *CommentVoteService) Create(newVote *models.CommentVote) (uint8, error) {
	signalForNotification := uint8(0)

	vote, err := s.repo.GetByUserId(newVote)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	// проверяем наличие vote
	if vote != 0 {
		signalForNotification += models.VoteSignalDelete
		
		err = s.repo.DeleteByUserId(newVote)
		if err != nil {
			return 0, err
		}
	}
	if vote != newVote.Vote {
		signalForNotification += models.VoteSignalCreate

		err = s.repo.Create(newVote)
		if err != nil {
			return 0, err
		}
	}
	return signalForNotification, nil
}
