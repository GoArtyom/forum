package handler

import (
	"fmt"
	"log"
	"net/http"

	"forum/internal/models"
	"forum/pkg"
)

// GET
func (h Handler) signinGET(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signin" {
		log.Printf("signin: not found %s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}
	if r.Method != http.MethodGet {
		log.Printf("signin: method not allowed %s\n", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed) // 405
		return
	}
	err := h.template.ExecuteTemplate(w, "index.html", fmt.Sprintf("Path:%s\nMethod:%s", r.URL.Path, r.Method))
	if err != nil {
		log.Printf("signin: execute %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
	}
}

// POST
func (h Handler) signinPOST(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/signin" {
		log.Printf("signinPost: not found %s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}
	if r.Method != http.MethodPost {
		log.Printf("signinPost: method not allowed %s\n", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed) // 405
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Printf("signinPost: parse form %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}
	// validate name/ email/ password
	user := &models.SignInUser{
		Email:    r.Form.Get("email"),
		Password: r.Form.Get("password"),
	}
	userId, err := h.service.SignInUser(user)
	if err != nil {
		log.Printf("signinPost: sign in user %s\n", err.Error())
		if err == models.IncorData {
			
			http.Redirect(w, r, "/signin", http.StatusSeeOther) // 303
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}

	session, err := h.service.CreateSession(userId)
	if err != nil {
		log.Printf("signinPost: create session %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}

	pkg.SetCookie(w, session.UUID, session.ExpireAt)

	http.Redirect(w, r, "/", http.StatusSeeOther) // 303
}
