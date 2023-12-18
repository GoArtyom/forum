package form

import (
	"fmt"
	"forum/internal/models"
	"log"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Form struct {
	Errors     map[string][]string
	DataForErr *models.DataForErr
	Request    *http.Request
}

func New(r *http.Request) *Form {
	return &Form{
		Errors: map[string][]string{},
		DataForErr: &models.DataForErr{
			Title:    r.Form.Get("title"),
			Content:  r.Form.Get("content"),
			Name:     r.Form.Get("name"),
			Email:    r.Form.Get("email"),
			Password: r.Form.Get("password"),
		},
		Request: r,
	}
}

func (f *Form) ErrLengthMax(key string, length int) {
	value := f.Request.Form.Get(key)
	if value == "" {
		return
	}

	if utf8.RuneCountInString(value) > length {
		err := fmt.Sprintf("This fielf is too long. Maximum is %d characters.", length)
		f.Errors[key] = append(f.Errors[key], err)
	}
}

func (f *Form) ErrLengthMin(key string, length int) {
	value := f.Request.Form.Get(key)
	if utf8.RuneCountInString(value) < length {
		err := fmt.Sprintf("This field is too short. Minimum is %d characters.", length)
		f.Errors[key] = append(f.Errors[key], err)
	}
}

func (f *Form) ErrEmpty(keys ...string) {
	for _, key := range keys {
		value := f.Request.Form.Get(key)
		value = strings.TrimSpace(value)
		if len(value) == 0 {
			f.Errors[key] = append(f.Errors[key], "This field is empty.")
		}
	}
}

func (f *Form) ErrLog(s string) {
	for key, errors := range f.Errors {
		for _, err := range errors {
			log.Printf(`%sKey="%s":%s`, s, key, err)
		}
	}
}

func (f *Form) ValidEmail(key string) {
	value := f.Request.Form.Get(key)
	p := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !p.MatchString(value) {
		f.Errors[key] = append(f.Errors[key], models.ErrEmail)
	}
}

func (f *Form) ValidPassword(key string) {
	value := f.Request.Form.Get(key)
	patterns := map[string]string{
		`[!@#$%^&*]`: "Contains at least one of the following special characters: !@#$%^&*",
		`\d`:         "Contains at least one number.",
		`[A-Z]`:      "Contains at least one uppercase letter.",
		`[a-z]`:      "Contains at least one lowercase letter.",
	}

	for pattern, err := range patterns {
		p := regexp.MustCompile(pattern)
		if !p.MatchString(value) {
			f.Errors[key] = append(f.Errors[key], err)
		}
	}
}
