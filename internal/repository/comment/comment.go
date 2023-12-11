package comment

import (
	"database/sql"

	"forum/internal/models"
)

type CommentSqlite struct {
	db *sql.DB
}

func NewCommentSqlite(db *sql.DB) *CommentSqlite {
	return &CommentSqlite{db: db}
}

func (r *CommentSqlite) CreateComment(comment *models.CreateComment) error {
	query := "SELECT id FROM posts WHERE id = $1"
	err := r.db.QueryRow(query, comment.PostId).Scan(&comment.PostId)
	if err != nil {
		return err
	}
	query2 := "INSERT INTO comments (post_id, content, user_id, user_name, create_at) VALUES($1, $2, $3, $4, $5)"
	_, err = r.db.Exec(query2, comment.PostId, comment.Content, comment.UserId, comment.UserName, comment.CreateAt)
	return err
}

func (r *CommentSqlite) GetAllCommentByPostId(postId int) ([]*models.Comment, error) {
	query := "SELECT * FROM comments WHERE post_id = $1"
	rows, err := r.db.Query(query, postId)
	if err != nil {
		return nil, err
	}
	comments := make([]*models.Comment, 0)
	for rows.Next() {
		comment := new(models.Comment)
		err := rows.Scan(&comment.Id, &comment.PostId, &comment.Content,
			&comment.UserId, &comment.UserName, &comment.CreateAt)
		if err != nil {
			return nil, err
		}
		// like & dislike
		query = "SELECT COALESCE(SUM(CASE WHEN vote = 1 THEN 1 ELSE 0 END), 0), COALESCE(SUM(CASE WHEN vote = -1 THEN 1 ELSE 0 END), 0) FROM comments_votes WHERE comment_id = $1"
		err = r.db.QueryRow(query, comment.Id).Scan(&comment.Like, &comment.Dislike)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
