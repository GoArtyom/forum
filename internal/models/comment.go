package models

import "time"

type CreateComment struct {
	PostId   int       `json:"post_id"`
	Content  string    `json:"content"`
	UserId   int       `json:"user_id"`
	UserName string    `json:"user_name"`
	CreateAt time.Time `json:"create_at"`
}

type Comment struct {
	Id       int       `json:"id"`
	PostId   int       `json:"post_id"`
	Content  string    `json:"content"`
	UserId   int       `json:"user_id"`
	UserName string    `json:"user_name"`
	Like     int       `json:"like"`
	Dislike  int       `json:"dislike"`
	CreateAt time.Time `json:"create_at"`
}

type CommentVote struct {
	CommentId int `json:"comment_id"`
	UserId    int `json:"user_id"`
	Vote      int `json:"vote"`
}
