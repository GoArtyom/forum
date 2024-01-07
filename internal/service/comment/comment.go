package comment

import (
	"forum/internal/models"
	repo "forum/internal/repository"
)

type CommentService struct {
	repo repo.Comment
}

func NewCommentService(repo repo.Comment) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) Create(comment *models.CreateComment) (int, error) {
	return s.repo.Create(comment)
}

func (s *CommentService) GetById(commentId int) (*models.Comment, error) {
	return s.repo.GetById(commentId)
}

func (s *CommentService) GetAllByPostId(postId int) ([]*models.Comment, error) {
	comments, err := s.repo.GetAllByPostId(postId)
	if err != nil {
		return nil, err
	}

	for i, j := 0, len(comments)-1; i < j; i, j = i+1, j-1 {
		comments[i], comments[j] = comments[j], comments[i]
	}

	return comments, nil
}

func (s *CommentService) GetAllByUserVote(userId int) ([]*models.Comment, error) {
	likeComments, err := s.repo.GetByVote(userId, models.VoteLike)
	if err != nil {
		return nil, err
	}

	dislikeComments, err := s.repo.GetByVote(userId, models.VoteDislike)
	if err != nil {
		return nil, err
	}

	for _, comment := range likeComments {
		comment.Like = models.VoteLike
	}

	for _, comment := range dislikeComments {
		comment.Like = models.VoteDislike
	}

	return append(likeComments, dislikeComments...), nil
}

func (s *CommentService) GetAllByUserId(userId int) ([]*models.Comment, error) {
	return s.repo.GetAllByUserId(userId)
}

func (s *CommentService) UpdateById(upComment *models.UpdateComment) error {
	comment, err := s.repo.GetById(upComment.Id)
	if err != nil {
		return err
	}

	if comment.UserId == upComment.UserId {
		return s.repo.UpdateById(upComment)
	}

	return models.ErrComment
}

func (s *CommentService) DeleteById(formDel *models.DeleteComment) error {
	if formDel.UserRole == models.AdminRole {
		return s.repo.DeleteById(formDel.CommentId)
	}

	comment, err := s.repo.GetById(formDel.CommentId)
	if err != nil {
		return err
	}

	if comment.UserId == formDel.UserId {
		return s.repo.DeleteById(formDel.CommentId)
	}

	return models.ErrComment
}
