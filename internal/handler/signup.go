package handler

import (
	"log"
	"net/http"

	"forum/internal/models"
)

func (h Handler) signup(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		log.Printf("signup: not found %s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}
	if r.Method != http.MethodGet {
		log.Printf("signup: method not allowed %s\n", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed) // 405
		return
	}
	err := h.template.ExecuteTemplate(w, "signup.html", nil)
	if err != nil {
		log.Printf("signup: execute %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
	}
}

func (h Handler) signupPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/signup" {
		log.Printf("signupPost: not found %s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}
	if r.Method != http.MethodPost {
		log.Printf("signupPost: method not allowed %s\n", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed) // 405
		return
	}
	if err := r.ParseForm(); err != nil {

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}
	// validate name/ email/ password
	user := &models.CreateUser{
		Name:     r.Form.Get("name"),
		Email:    r.Form.Get("email"),
		Password: r.Form.Get("password"),
	}
	err := h.service.CreateUser(user)
	if err != nil {
		// validate err
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}
	http.Redirect(w, r, "/signin", http.StatusSeeOther) // 303
}
