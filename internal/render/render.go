package render

import (
	"forum/internal/models"
	"forum/pkg/form"
)

type Data struct {
	User       *models.User       `json:"user"`
	Post       *models.Post       `json:"post"`
	Posts      []*models.Post     `json:"posts"`
	Comments   []*models.Comment  `json:"coments"`
	Categories []*models.Category `json:"categories"`
	Form       *form.Form         `json:"form"`
}
