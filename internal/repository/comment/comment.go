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

func (r *CommentSqlite) Create(comment *models.CreateComment) (int, error) {

	query2 := "INSERT INTO comments (post_id, content, user_id, user_name, create_at) VALUES($1, $2, $3, $4, $5)"
	res, err := r.db.Exec(query2, comment.PostId, comment.Content, comment.UserId, comment.UserName, comment.CreateAt)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *CommentSqlite) GetById(commentId int) (*models.Comment, error) {
	comment := new(models.Comment)

	query := "SELECT * FROM comments WHERE id = ?"

	err := r.db.QueryRow(query, commentId).Scan(&comment.Id, &comment.PostId, &comment.Content,
		&comment.UserId, &comment.UserName, &comment.CreateAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrComment
		}
		return nil, err
	}

	// like & dislike
	query = "SELECT COALESCE(SUM(CASE WHEN vote = 1 THEN 1 ELSE 0 END), 0), COALESCE(SUM(CASE WHEN vote = -1 THEN 1 ELSE 0 END), 0) " +
		"FROM comments_votes WHERE comment_id = $1"

	err = r.db.QueryRow(query, comment.Id).Scan(&comment.Like, &comment.Dislike)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (r *CommentSqlite) GetByVote(userId, vote int) ([]*models.Comment, error) {
	query := "SELECT c.id, c.post_id, c.content, c.user_id, c.user_name, c.create_at " +
		"FROM comments c JOIN comments_votes cv ON c.id = cv.comment_id  WHERE cv.user_id = $1 AND cv.vote = $2"

	rows, err := r.db.Query(query, userId, vote)
	if err != nil {
		return nil, err
	}

	comments := make([]*models.Comment, 0)

	for rows.Next() {
		comment := new(models.Comment)

		err := rows.Scan(&comment.Id, &comment.PostId, &comment.Content, &comment.UserId,
			&comment.UserName, &comment.CreateAt)
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

func (r *CommentSqlite) GetAllByPostId(postId int) ([]*models.Comment, error) {
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
		query = "SELECT COALESCE(SUM(CASE WHEN vote = 1 THEN 1 ELSE 0 END), 0), COALESCE(SUM(CASE WHEN vote = -1 THEN 1 ELSE 0 END), 0) " +
			"FROM comments_votes WHERE comment_id = $1"

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

func (r *CommentSqlite) GetAllByUserId(userId int) ([]*models.Comment, error) {
	query := "SELECT * FROM comments WHERE user_id = $1"

	rows, err := r.db.Query(query, userId)
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

		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *CommentSqlite) DeleteById(commentId int) error {
	query := "DELETE FROM comments WHERE id = ?"

	res, err := r.db.Exec(query, commentId)
	if err != nil {
		return err
	}

	countDel, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if countDel == 0 {
		return models.ErrComment
	}

	return nil
}

func (r *CommentSqlite) UpdateById(upComment *models.UpdateComment) error {
	query := "UPDATE comments SET content = $1 WHERE id = $2"

	res, err := r.db.Exec(query, upComment.Content, upComment.Id)
	if err != nil {
		return err
	}

	countUpd, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if countUpd == 0 {
		return models.ErrComment
	}

	return nil
}
