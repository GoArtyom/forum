package postvote

import (
	"database/sql"

	"forum/internal/models"
	"forum/internal/repository"
)

type PostVoteService struct {
	repo repository.PostVote
}

func NewPostVoteService(repo repository.PostVote) *PostVoteService {
	return &PostVoteService{repo: repo}
}

func (s *PostVoteService) CreatePostVote(newVote *models.PostVote) error {
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
		err = s.repo.CreatePostVote(newVote)
		if err != nil {
			return err
		}
	}
	return nil
}
