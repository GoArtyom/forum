package session

import (
	"time"

	"forum/internal/models"
	repo "forum/internal/repository"

	"github.com/gofrs/uuid"
)

type SessionService struct {
	repo repo.Session
}

func NewSessionService(repo repo.Session) *SessionService {
	return &SessionService{repo: repo}
}

func (s *SessionService) Create(userId int) (*models.Session, error) {
	oldSession, _ := s.repo.GetByUserId(userId)
	if oldSession != nil {
		err := s.repo.DeleteByUUID(oldSession.UUID)
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
		ExpireAt: time.Now().Add(time.Hour * 24),
	}

	err = s.repo.Create(newSession)
	if err != nil {
		return nil, err
	}

	return newSession, nil
}

func (s *SessionService) GetByUUID(uuid string) (*models.Session, error) {
	return s.repo.GetByUUID(uuid)
}

func (s *SessionService) DeleteByUUID(uuid string) error {
	return s.repo.DeleteByUUID(uuid)
}
