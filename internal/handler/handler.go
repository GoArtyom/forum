package handler

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"forum/internal/models"
	"forum/internal/service"
)

type Handler struct {
	service  *service.Service
	template *template.Template
}

func NewHandler(service *service.Service, tpl *template.Template) *Handler {
	return &Handler{
		service:  service,
		template: tpl,
	}
}

func (h *Handler) getUserFromContext(r *http.Request) *models.User {
	user, ok := r.Context().Value(keyUser).(*models.User)
	if !ok {
		return nil
	}
	return user
}

func (h *Handler) getPostIdFromURL(path string) (int, error) {
	parts := strings.Split(path, "/")
	if len(parts) != 3 {
		return 0, errors.New("incorrect path")
	}
	postId, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, err
	}
	return postId, nil
}
