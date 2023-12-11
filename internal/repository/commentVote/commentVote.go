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

func (r *CommentVoteSqlite) CreateCommentVote(newVote *models.CommentVote) error {
	query := "INSERT INTO comments_votes (comment_id, user_id, vote) VALUES($1, $2, $3)"
	_, err := r.db.Exec(query, newVote.CommentId, newVote.UserId, newVote.Vote)
	return err
}

func (r *CommentVoteSqlite) GetVoteByUserId(newVote *models.CommentVote) (int, error) {
	var vote int

	query := "SELECT vote FROM comments_votes WHERE comment_id = $1 AND user_id = $2"
	err := r.db.QueryRow(query, newVote.CommentId, newVote.UserId).Scan(&vote)
	return vote, err
}

func (r *CommentVoteSqlite) DeleteVoteByUserId(newVote *models.CommentVote) error {
	query := "DELETE FROM comments_votes WHERE comment_id = $1 AND user_id = $2"
	_, err := r.db.Exec(query, newVote.CommentId, newVote.UserId)
	return err
}
