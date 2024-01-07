package render

import (
	"forum/internal/models"
	"forum/pkg/form"
)

type MainData struct {
	User       *models.User       `json:"user"`
	Posts      []*models.Post     `json:"posts"`
	Categories []*models.Category `json:"categories"`
}

type CreatePost struct {
	User       *models.User       `json:"user"`
	Categories []*models.Category `json:"categories"`
	Form       *form.Form         `json:"form"`
}

type UpdatePost struct {
	User       *models.User       `json:"user"`
	Post       *models.Post       `json:"post"`
	Categories []*models.Category `json:"categories"`
	Form       *form.Form         `json:"form"`
}
type UpdateComment struct {
	User    *models.User    `json:"user"`
	Comment *models.Comment `json:"comment"`
	Form    *form.Form      `json:"form"`
}

type OnePostData struct {
	User     *models.User      `json:"user"`
	Post     *models.Post      `json:"post"`
	Comments []*models.Comment `json:"coments"`
	Form     *form.Form        `json:"form"`
}

type OnlyForm struct {
	User *models.User `json:"user"`
	Form *form.Form   `json:"form"`
}

type MyActivity struct {
	User          *models.User      `json:"user"`
	MyPosts       []*models.Post    `json:"posts"`
	MyComments    []*models.Comment `json:"comments"`
	VotesPosts    []*models.Post    `json:"vote_posts"`
	VotesComments []*models.Comment `json:"vote_comments"`
}

type AdminData struct {
	User             *models.User       `json:"user"`
	Reports          []*models.Report   `json:"reports"`
	Users            []*models.User     `json:"users"`
	RequestModerator []*models.User     `json:"request_moderator"`
	Categories       []*models.Category `json:"categories"`
	Form             *form.Form         `json:"form"`
}

type MyNotifications struct {
	User          *models.User           `json:"user"`
	Notifications []*models.Notification `json:"notifications"`
}
