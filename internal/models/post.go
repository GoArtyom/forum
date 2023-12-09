package models

import "time"

type CreatePost struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	UserId     int    `json:"user_id"`
	UserName   string `json:"user_name"`
	Categories *[]string
	CreateAt   time.Time `json:"create_at"`
}

type Post struct {
	PostId   int       `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	UserId   int       `json:"user_id"`
	UserName string    `json:"user_name"`
	CreateAt time.Time `json:"create_at"`
}
