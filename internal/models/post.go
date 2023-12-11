package models

import "time"

type CreatePost struct {
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	UserId     int       `json:"user_id"`
	UserName   string    `json:"user_name"`
	Categories []string  `json:"categories"`
	CreateAt   time.Time `json:"create_at"`
}

type Post struct {
	PostId     int        `json:"id"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	UserId     int        `json:"user_id"`
	UserName   string     `json:"user_name"`
	CreateAt   time.Time  `json:"create_at"`
	Categories []string   `json:"categories"`
	Comments   []*Comment `json:"comments"`
	Like       int        `json:"like"`
	Dislike    int        `json:"dislike"`
}

type PostVote struct {
	PostId int `json:"post_id"`
	UserId int `json:"user_id"`
	Vote   int `json:"vote"`
}

type Data struct {
	User       *User       `json:"user"`
	Post       *Post       `json:"post"`
	Posts      []*Post     `json:"posts"`
	Comments   []*Comment  `json:"coments"`
	Categories []*Category `json:"categories"`
}
