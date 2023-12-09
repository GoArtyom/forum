package handler

import (
	"html/template"
	"net/http"

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
