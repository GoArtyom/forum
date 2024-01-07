package models

import (
	"time"
)

const (
	NoticeTypeComment     = uint8(1)
	NoticeTypePostVote    = uint8(2)
	NoticeTypeCommentVote = uint8(3)
)
const (
	VoteSignalDelete = uint8(1)
	VoteSignalCreate = uint8(2)
	VoteSignalChange = uint8(3)
)

const (
	DelNoticById        = uint8(1)
	DelNoticByAuthorAll = uint8(2)
	DelNoticByUser      = uint8(3)
)

type Notification struct {
	Id        int       `json:"id"`
	PostId    int       `json:"post_id"`
	CommentId int       `json:"comment_id"`
	AuthorId  int       `json:"author_id"`
	UserId    int       `json:"user_id"`
	UserName  string    `json:"user_name"`
	Content   string    `json:"content"` 
	Vote      int       `json:"vote"`
	Type      uint8     `json:"type"`
	CreateAt  time.Time `json:"create_at"`
}

type DeleteNotification struct {
	Id     int   `json:"id"`
	PostId int   `json:"post_id"`
	UserId int   `json:"user_ids"`
	Type   uint8 `json:"type"`
	Method uint8 `json:"method"`
}
