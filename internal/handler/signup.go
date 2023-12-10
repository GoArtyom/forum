package handler

import (
	"fmt"
	"log"
	"net/http"

	"forum/internal/models"
)

// GET
func (h *Handler) signupGET(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		log.Printf("signupGET: not found %s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}
	if r.Method != http.MethodGet {
		log.Printf("signupGET: method not allowed %s\n", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed) // 405
		return
	}
	err := h.template.ExecuteTemplate(w, "index.html", fmt.Sprintf("Path:%s\nMethod:%s", r.URL.Path, r.Method))
	if err != nil {
		log.Printf("signupGET: ExecuteTemplate %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
	}
}

// POST
func (h *Handler) signupPOST(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/signup" {
		log.Printf("signupPOST: not found %s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}
	if r.Method != http.MethodPost {
		log.Printf("signupPOST: method not allowed %s\n", r.Method)
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
		log.Printf("signupPOST: create user: %s\n", err.Error())
		if err.Error() == models.UniqueUser.Error() {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}
	http.Redirect(w, r, "/signin", http.StatusSeeOther) // 303
}
