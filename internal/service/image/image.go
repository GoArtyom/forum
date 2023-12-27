package image

import (
	"forum/internal/models"
	repo "forum/internal/repository"
	"io"
	"os"
	"strings"

	"github.com/gofrs/uuid"
)

type ImageService struct {
	repo repo.Image
}

func NewImageService(repo repo.Image) *ImageService {
	return &ImageService{repo: repo}
}

func (s *ImageService) CreateImageByPostIt(newImage *models.CreateImage) error {
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

	err = s.repo.CreateImageByPostId(newImage)
	if err != nil {
		return err
	}

	newFile, err := os.OpenFile("ui/static/img/"+newImage.Name+"."+newImage.Type, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		err = s.repo.DeleteImageByPostId(newImage.PostId)
		if err != nil {
			return err
		}
		return err
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, content)
	if err != nil {
		s.repo.DeleteImageByPostId(newImage.PostId)
		if err != nil {
			return err
		}
		return err
	}
	return nil
}

func (s *ImageService) GetImageByPostId(postId int) (*models.Image, error) {
	return s.repo.GetImageByPostId(postId)
}
