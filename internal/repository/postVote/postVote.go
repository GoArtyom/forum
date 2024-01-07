package postvote

import (
	"database/sql"

	"forum/internal/models"
)

type PostVoteSqlite struct {
	db *sql.DB
}

func NewPostVoteSqlite(db *sql.DB) *PostVoteSqlite {
	return &PostVoteSqlite{db: db}
}

func (r *PostVoteSqlite) Create(newVote *models.PostVote) error {
	query := "INSERT INTO posts_votes (post_id, user_id, vote) VALUES($1, $2, $3)"
	_, err := r.db.Exec(query, newVote.PostId, newVote.UserId, newVote.Vote)
	return err
}

func (r *PostVoteSqlite) GetByUserId(newVote *models.PostVote) (int, error) {
	var vote int

	query := "SELECT vote FROM posts_votes WHERE post_id = $1 AND user_id = $2"
	err := r.db.QueryRow(query, newVote.PostId, newVote.UserId).Scan(&vote)
	return vote, err
}

func (r *PostVoteSqlite) DeleteByUserId(newVote *models.PostVote) error {
	query := "DELETE FROM posts_votes WHERE post_id = $1 AND user_id = $2"
	_, err := r.db.Exec(query, newVote.PostId, newVote.UserId)
	return err
}
