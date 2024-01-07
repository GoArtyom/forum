package form

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"

	"forum/internal/models"
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
			Category: r.Form.Get("category"),
			Name:     r.Form.Get("name"),
			Email:    r.Form.Get("email"),
			Password: r.Form.Get("password"),
			Report: r.Form.Get("report"),
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
		err := fmt.Sprintf(models.TextErrFieldTooLong, length)
		f.Errors[key] = append(f.Errors[key], err)
	}
}

func (f *Form) ErrLengthMin(key string, length int) {
	value := f.Request.Form.Get(key)
	if utf8.RuneCountInString(value) < length {
		err := fmt.Sprintf(models.TextErrFieldTooShort, length)
		f.Errors[key] = append(f.Errors[key], err)
	}
}

func (f *Form) ErrEmpty(keys ...string) {
	for _, key := range keys {
		value := f.Request.Form.Get(key)
		valueWithoutSpace := strings.TrimSpace(value)
		if len(valueWithoutSpace) == 0 {
			f.Errors[key] = append(f.Errors[key], models.TextErrEmptyField)
		} else if value != valueWithoutSpace {
			f.Errors[key] = append(f.Errors[key], models.TextErrExtraSpace)
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
		f.Errors[key] = append(f.Errors[key], models.TextErrEmail)
	}
}

func (f *Form) ValidPassword(key string) {
	value := f.Request.Form.Get(key)
	patterns := map[string]string{
		`[!@#$%^&*]`: models.TextErrSpecialChars,
		`\d`:         models.TextErrOneNumber,
		`[A-Z]`:      models.TextErrUpLetter,
		`[a-z]`:      models.TextErrLowLetter,
	}

	for pattern, err := range patterns {
		p := regexp.MustCompile(pattern)
		if !p.MatchString(value) {
			f.Errors[key] = append(f.Errors[key], err)
		}
	}
}

func (f *Form) ErrImg(h *multipart.FileHeader) {
	if h.Size > 20<<20 {
		f.Errors["img"] = append(f.Errors["img"], models.TextErrSize)
	}
	nameSplit := strings.Split(h.Filename, ".")
	mime := strings.ToLower(nameSplit[len(nameSplit)-1])
	if len(nameSplit) != 2 || (mime != "gif" && mime != "png" && mime != "jpg" && mime != "jpeg") {
		f.Errors["img"] = append(f.Errors["img"], fmt.Sprintf(models.TextErrTypeFile, mime))
	}
}
