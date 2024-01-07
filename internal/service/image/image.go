package image

import (
	"io"
	"os"
	"strings"

	"forum/internal/models"
	repo "forum/internal/repository"

	"github.com/gofrs/uuid"
)

type ImageService struct {
	repo repo.Image
}

func NewImageService(repo repo.Image) *ImageService {
	return &ImageService{repo: repo}
}

func (s *ImageService) CreateByPostId(newImage *models.CreateImage) error {
	splitName := strings.Split(newImage.Header.Filename, ".")

	hashName, err := uuid.NewV4()
	if err != nil {
		return err
	}

	newImage.Name = hashName.String()
	newImage.Type = strings.ToLower(splitName[1])

	content, err := newImage.Header.Open()
	if err != nil {
		return err
	}
	defer content.Close()

	err = s.repo.CreateByPostId(newImage)
	if err != nil {
		return err
	}

	newFile, err := os.OpenFile("ui/static/img/"+newImage.Name+"."+newImage.Type, os.O_WRONLY|os.O_CREATE, 0o666)
	if err != nil {
		err = s.repo.DeleteByPostId(newImage.PostId)
		if err != nil {
			return err
		}
		return err
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, content)
	if err != nil {
		s.repo.DeleteByPostId(newImage.PostId)
		if err != nil {
			return err
		}
		return err
	}
	return nil
}

func (s *ImageService) GetByPostId(postId int) (*models.Image, error) {
	return s.repo.GetByPostId(postId)
}

func (s *ImageService) DeleteByPostId(PostId int) error {
	image, err := s.repo.GetByPostId(PostId)
	if err != nil {
		return err
	}
	err = os.Remove("ui/static/img/" + image.Name + "." + image.Type)
	if err != nil {
		return err
	}

	return s.repo.DeleteByPostId(PostId)
}
