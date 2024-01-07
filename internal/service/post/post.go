package post

import (
	"database/sql"
	"forum/internal/models"
	repo "forum/internal/repository"
	"os"
)

type PostService struct {
	post repo.Post
	cat  repo.Category
	img  repo.Image
}

func NewPostService(repo *repo.Repository) *PostService {
	return &PostService{
		post: repo.Post,
		cat:  repo.Category,
		img:  repo.Image,
	}
}

func (s *PostService) Create(post *models.CreatePost) (int, error) {
	return s.post.Create(post)
}

func (s *PostService) GetById(postId int) (*models.Post, error) {
	return s.post.GetById(postId)
}

func (s *PostService) GetAll() ([]*models.Post, error) {
	return s.post.GetAll()
}

func (s *PostService) GetAllByUserId(userId int) ([]*models.Post, error) {
	return s.post.GetAllByUserId(userId)
}

func (s *PostService) GetAllByUserVote(userId int) ([]*models.Post, error) {
	likePosts, err := s.post.GetByVote(userId, models.VoteLike)
	if err != nil {
		return nil, err
	}

	dislikePosts, err := s.post.GetByVote(userId, models.VoteDislike)
	if err != nil {
		return nil, err
	}

	for _, post := range likePosts {
		post.Like = models.VoteLike
	}

	for _, post := range dislikePosts {
		post.Like = models.VoteDislike
	}

	return append(likePosts, dislikePosts...), nil
}

func (s *PostService) GetByCategory(category string) ([]*models.Post, error) {
	_, err := s.cat.GetByName(category)
	if err != nil {
		return nil, err
	}

	return s.post.GetByCategory(category)
}

func (s *PostService) GetByVote(userId, vote int) ([]*models.Post, error) {
	return s.post.GetByVote(userId, vote)
}

func (s *PostService) UpdateById(upPost *models.UpdatePost) error {
	post, err := s.post.GetById(upPost.PostId)
	if err != nil {
		return err
	}
	if post.UserId == upPost.UserId {
		return s.post.UpdateById(upPost)
	}
	return models.ErrPost
}

func (s *PostService) DeleteById(formDel *models.DeletePost) error {
	img, err := s.img.GetByPostId(formDel.PostId)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if formDel.ServerErr {
		if img != nil {
			err = os.Remove("ui/static/img/" + img.Name + "." + img.Type)
			if err != nil {
				return err
			}
		}
		return s.post.DeleteById(formDel.PostId)
	}
	if formDel.UserRole >= models.ModeratorRole {
		if img != nil {
			err = os.Remove("ui/static/img/" + img.Name + "." + img.Type)
			if err != nil {
				return err
			}
		}
		return s.post.DeleteById(formDel.PostId)
	}

	post, err := s.post.GetById(formDel.PostId)
	if err != nil {
		return err
	}

	if post.UserId == formDel.UserId {
		if img != nil {
			err = os.Remove("ui/static/img/" + img.Name + "." + img.Type)
			if err != nil {
				return err
			}
		}
		return s.post.DeleteById(formDel.PostId)
	}
	return models.ErrPost
}
