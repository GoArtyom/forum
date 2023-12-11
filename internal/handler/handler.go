package handler

import (
	"errors"
	"fmt"
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

func (h *Handler) getVote(voteStr string) (int, error) {
	vote, err := strconv.Atoi(voteStr)
	if err != nil {
		return 0, err
	}
	if vote != 1 && vote != -1 {
		return 0, fmt.Errorf("incorrect request vote = %d", vote)
	}
	return vote, nil
}

func (h *Handler) getPostIdFromForm(postIdStr string) (int, error) {
	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		return 0, err
	}
	if postId < 1 {
		return 0, fmt.Errorf("incorrect request postId = %d", postId)
	}
	return postId, nil
}
