package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"forum/internal/models"
)

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
	if postId < 1 {
		return 0, errors.New("post id is less than 1")
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

func (h *Handler) getIntFromForm(r *http.Request, key string) (int, error) {
	value := r.Form.Get(key)
	if value == "" {
		return -1, fmt.Errorf(`empty value="%s"`, key)
	}
	id, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	if id < 1 {
		return 0, fmt.Errorf(`in the variable %[1]s the value is less than one. "%[1]s"=%[2]d`, key, id)
	}
	return id, nil
}
