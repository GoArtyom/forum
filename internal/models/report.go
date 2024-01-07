package models

import "time"

const (
	Approve = 1
	Refuse  = 2
)

type CreateReport struct {
	PostId        int       `json:"post_id"`
	Content       string    `json:"content"`
	ModeratorId   int       `json:"moderator_id"`
	ModeratorName string    `json:"moderator_name"`
	CreateAt      time.Time `json:"create_at"`
}

type Report struct {
	Id            int       `json:"id"`
	PostId        int       `json:"post_id"`
	Content       string    `json:"content"`
	ModeratorId   int       `json:"moderator_id"`
	ModeratorName string    `json:"moderator_name"`
	CreateAt      time.Time `json:"create_at"`
}
