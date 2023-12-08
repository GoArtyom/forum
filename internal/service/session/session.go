package session

import (
	"time"

	"forum/internal/models"
	"forum/internal/repository"

	"github.com/gofrs/uuid"
)

type SessionServise struct {
	repo repository.Session
}

func NewSessionServise(repo repository.Session) *SessionServise {
	return &SessionServise{repo: repo}
}

func (s *SessionServise) CreateSession(userId int) (*models.Session, error) {
	oldSession, _ := s.repo.GetSessionByUserId(userId)
	if oldSession != nil {
		err := s.repo.DeleteSessionByUUID(oldSession.UUID)
		if err != nil {
			return nil, err
		}
	}

	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	newSession := &models.Session{
		User_id:  userId,
		UUID:     uuid.String(),
		ExpireAt: time.Now().Add(time.Hour* 24),
	}

	err = s.repo.CreateSession(newSession)
	if err != nil {
		return nil, err
	}

	return newSession, nil
}
