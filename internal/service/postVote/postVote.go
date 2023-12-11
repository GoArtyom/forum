package postvote

import (
	"database/sql"

	"forum/internal/models"
	repo "forum/internal/repository"
)

type PostVoteService struct {
	repo repo.PostVote
}

func NewPostVoteService(repo repo.PostVote) *PostVoteService {
	return &PostVoteService{repo: repo}
}

func (s *PostVoteService) CreatePostVote(newVote *models.PostVote) error {
	vote, err := s.repo.GetVoteByUserIdR(newVote)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if vote != 0 { // проверяем наличие vote
		err = s.repo.DeleteVoteByUserIdR(newVote)
		if err != nil {
			return err
		}
	}
	if vote != newVote.Vote {
		err = s.repo.CreatePostVote(newVote)
		if err != nil {
			return err
		}
	}
	return nil
}
