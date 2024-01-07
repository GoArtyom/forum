package notification

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type NoticService struct {
	notic   repository.Notification
	post    repository.Post
	comment repository.Comment
}

func NewNoticService(repo *repository.Repository) *NoticService {
	return &NoticService{
		notic:   repo.Notification,
		post:    repo.Post,
		comment: repo.Comment,
	}
}

func (s *NoticService) Create(notic *models.Notification) error {
	switch notic.Type {
	case models.NoticeTypeComment:
		post, err := s.post.GetById(notic.PostId)
		if err != nil {
			return err
		}

		if post.UserId != notic.UserId {
			notic.AuthorId = post.UserId
			return s.notic.Create(notic)
		}

	case models.NoticeTypePostVote:
		post, err := s.post.GetById(notic.PostId)
		if err != nil {
			return err
		}

		if post.UserId != notic.UserId {
			notic.AuthorId = post.UserId
			notic.Content = post.Title

			return s.notic.Create(notic)
		}

	case models.NoticeTypeCommentVote:
		comment, err := s.comment.GetById(notic.CommentId)
		if err != nil {
			return err
		}
		if comment.UserId != notic.UserId {
			notic.AuthorId = comment.UserId
			notic.Content = comment.Content

			return s.notic.Create(notic)
		}
	}
	return nil
}

func (s *NoticService) GetAllByAuthorId(authorId int) ([]*models.Notification, error) {
	return s.notic.GetAllByAuthorId(authorId)
}

func (s *NoticService) GetCountByAuthorId(authorId int) (int, error) {
	return s.notic.GetCountByAuthorId(authorId)
}

func (s *NoticService) Delete(formDel *models.DeleteNotification) error {
	switch formDel.Method {
	case models.DelNoticById:
		return s.notic.DeleteById(formDel.Id)
	case models.DelNoticByAuthorAll:
		return s.notic.DeleteByAuthorId(formDel.UserId)
	case models.DelNoticByUser:
		return s.notic.DeleteByUserId(formDel)
	}
	return models.ErrNotification
}
