package session

import (
	"database/sql"

	"forum/internal/models"
)

type SessionSqlite struct {
	db *sql.DB
}

func NewSessionSqlite(db *sql.DB) *SessionSqlite {
	return &SessionSqlite{db: db}
}

func (r *SessionSqlite) CreateSession(session *models.Session) error {
	query := "INSERT INTO sessions (user_id, uuid, expire_at) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(query, session.User_id, session.UUID, session.ExpireAt)
	return err
}

func (r *SessionSqlite) DeleteSessionByUUID(sessionId string) error {
	query := "DELETE FROM sessions WHERE uuid = $1"
	_, err := r.db.Exec(query, sessionId)
	return err
}

func (r *SessionSqlite) GetSessionByUserId(userId int) (*models.Session, error) {
	session := models.Session{}
	query := "SELECT * FROM sessions WHERE user_id = $1"
	err := r.db.QueryRow(query, userId).Scan(&session.User_id, &session.UUID, &session.ExpireAt)
	return &session, err
}
