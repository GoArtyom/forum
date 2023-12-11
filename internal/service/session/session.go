package session

import (
	"time"

	"forum/internal/models"
	repo"forum/internal/repository"

	"github.com/gofrs/uuid"
)

type SessionService struct {
	repo repo.Session
}

func NewSessionService(repo repo.Session) *SessionService {
	return &SessionService{repo: repo}
}

func (s *SessionService) CreateSession(userId int) (*models.Session, error) {
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
		ExpireAt: time.Now().Add(time.Hour * 60),
	}

	err = s.repo.CreateSession(newSession)
	if err != nil {
		return nil, err
	}

	return newSession, nil
}

func (s *SessionService) GetSessionByUUID(uuid string) (*models.Session, error) {
	return s.repo.GetSessionByUUID(uuid)
}

func (s *SessionService) DeleteSessionByUUID(uuid string) error {
	return s.repo.DeleteSessionByUUID(uuid)
}
