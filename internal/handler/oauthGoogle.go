package handler

import (
	"fmt"
	"net/http"
)

const (
	googleOAuth = "https://accounts.google.com/o/oauth2/v2/auth"
)

func (h *Handler) signinGoogle(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s", googleOAuth)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) callbackGoogle(w http.ResponseWriter, r *http.Request) {
}
