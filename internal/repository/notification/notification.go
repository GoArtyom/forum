package notification

import (
	"database/sql"

	"forum/internal/models"
)

type NoticSqlite struct {
	db *sql.DB
}

func NewNoticService(db *sql.DB) *NoticSqlite {
	return &NoticSqlite{db: db}
}

func (r *NoticSqlite) Create(notic *models.Notification) error {
	var commentIdInterface interface{}
	if notic.CommentId == 0 {
		commentIdInterface = nil
	} else {
		commentIdInterface = notic.CommentId
	}

	query := "INSERT INTO notifications (post_id, comment_id, author_id, user_id, user_name, content, vote, type, create_at) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := r.db.Exec(query, notic.PostId, commentIdInterface, notic.AuthorId, notic.UserId,
		notic.UserName, notic.Content, notic.Vote, notic.Type, notic.CreateAt)
	return err
}

func (r *NoticSqlite) GetCountByAuthorId(authorId int) (int, error) {
	query := "SELECT COUNT(*) FROM notifications WHERE author_id = ?"

	var count int
	err := r.db.QueryRow(query, authorId).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *NoticSqlite) GetAllByAuthorId(authorId int) ([]*models.Notification, error) {
	query := "SELECT * FROM notifications WHERE author_id = ?"

	rows, err := r.db.Query(query, authorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notics []*models.Notification

	for rows.Next() {
		var commentIdInterface interface{}
		notic := &models.Notification{}

		err := rows.Scan(&notic.Id, &notic.PostId, &commentIdInterface, &notic.AuthorId, &notic.UserId,
			&notic.UserName, &notic.Content, &notic.Vote, &notic.Type, &notic.CreateAt)
		if err != nil {
			return nil, err
		}
		if commentIdInterface != nil {
			notic.CommentId = int(commentIdInterface.(int64))
		}

		notics = append(notics, notic)
	}

	return notics, nil
}

func (r *NoticSqlite) DeleteByAuthorId(authorId int) error {
	query := "DELETE FROM notifications WHERE author_id = ?"

	_, err := r.db.Exec(query, authorId)
	return err
}

func (r *NoticSqlite) DeleteById(noticId int) error {
	query := "DELETE FROM notifications WHERE id = ?"

	_, err := r.db.Exec(query, noticId)
	return err
}

func (r *NoticSqlite) DeleteByUserId(formDel *models.DeleteNotification) error {
	query := "DELETE FROM notifications WHERE post_id = ? AND user_id = ? AND type = ?"
	_, err := r.db.Exec(query, formDel.PostId, formDel.UserId, formDel.Type)
	return err
}
