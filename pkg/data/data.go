package data

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"

	"forum/internal/models"
)

type Data struct {
	User       *models.User        `json:"user"`
	Post       *models.Post        `json:"post"`
	Posts      []*models.Post      `json:"posts"`
	Comments   []*models.Comment   `json:"coments"`
	Categories []*models.Category  `json:"categories"`
	Errors     map[string][]string `json:"errors"`
}

func (d *Data) ErrLengthMax(r *http.Request, key string, length int) {
	value := r.Form.Get(key)
	if value == "" {
		return
	}

	if utf8.RuneCountInString(value) > length {
		err := fmt.Sprintf("This field is too long. Maximum is %d characters.", length)
		d.Errors[key] = append(d.Errors[key], err)
	}
}

func (d *Data) ErrLengthMin(r *http.Request, key string, length int) {
	value := r.Form.Get(key)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < length {
		err := fmt.Sprintf("This field is too short. Minimum is %d characters.", length)
		d.Errors[key] = append(d.Errors[key], err)
	}
}

func (d *Data) ErrEmpty(r *http.Request, keys ...string) {
	for _, key := range keys {
		value := r.Form.Get(key)
		value = strings.TrimSpace(value)
		if len(value) == 0 {
			d.Errors[key] = append(d.Errors[key], "This field is empty.")
		}
	}
}

func (d *Data) ErrLog(s string) {
	for key, errors := range d.Errors {
		for _, err := range errors {
			log.Printf(`%sKey="%s":%s`, s, key, err)
		}
	}
}



func (d *Data) IsValid(r *http.Request, key string, p *regexp.Regexp) {
	value := r.Form.Get(key)
	if value == "" {
		return
	}

	p.MatchString(value)
	if !p.MatchString(value) {
		switch key {
		case "email":
			d.Errors[key] = append(d.Errors[key], models.ErrEmail)
		case "password":
			d.Errors[key] = append(d.Errors[key], models.ErrPassword)
		}
	}
}
