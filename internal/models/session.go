package models

import "time"

type Session struct {
	User_id  int       `json:"user_id"`
	UUID     string    `json:"uuid"`
	ExpireAt time.Time `json:"expire_at"`
}
