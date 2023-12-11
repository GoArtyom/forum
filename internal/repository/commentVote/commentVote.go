package commentvote

import (
	"database/sql"

	"forum/internal/models"
)

type CommentVoteSqlite struct {
	db *sql.DB
}

func NewCommentVoteSqlite(db *sql.DB) *CommentVoteSqlite {
	return &CommentVoteSqlite{db: db}
}

func (s *CommentVoteSqlite) CreateCommentVote(newVote *models.CommentVote) error {
	return nil
}
