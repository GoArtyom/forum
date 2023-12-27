package models

import "mime/multipart"

type CreateImage struct {
	Header *multipart.FileHeader
	PostId int
	Name   string
	Type   string
}

type Image struct {
	Name string
	Type string
}
