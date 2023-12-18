package handler

import (
	"log"
	"net/http"

	"forum/internal/models"
	"forum/internal/render"
	"forum/pkg/form"
)

// GET
func (h *Handler) signupGET(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		log.Printf("signupGET:StatusNotFound%s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}
	if r.Method != http.MethodGet {
		log.Printf("signupGET:StatusMethodNotAllowed:%s\n", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed) // 405
		return
	}

	h.renderPage(w, "signup.html", &render.Data{})
}

// POST
func (h *Handler) signupPOST(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/signup" {
		log.Printf("signupPOST:StatusNotFound%s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}
	if r.Method != http.MethodPost {
		log.Printf("signupPOST:StatusMethodNotAllowed:%s\n", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed) // 405
		return
	}
	if err := r.ParseForm(); err != nil {

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}

	// validate name/ email/ password
	form := form.New(r)
	form.Errors = map[string][]string{}
	form.ErrEmpty("name", "email", "password")
	form.ErrLengthMin("name", 5)
	form.ErrLengthMax("name", 20)
	form.ErrLengthMin("email", 5)
	form.ErrLengthMax("email", 40)
	form.ErrLengthMin("password", 8)
	form.ErrLengthMax("password", 20)
	form.ValidEmail("email")
	form.ValidPassword("password")

	if len(form.Errors) != 0 {
		w.WriteHeader(http.StatusBadRequest) // 400
		form.ErrLog("signupPOST:")

		h.renderPage(w, "signup.html", &render.Data{
			Form: form,
		})
		return
	}

	user := &models.CreateUser{
		Name:     r.Form.Get("name"),
		Email:    r.Form.Get("email"),
		Password: r.Form.Get("password"),
	}
	err := h.service.CreateUser(user)
	if err != nil {
		switch err.Error() {
		case models.UniqueName:
			form.Errors["name"] = append(form.Errors["name"], "The user with that name has already been registered.")
		case models.UniqueEmail:
			form.Errors["email"] = append(form.Errors["email"], "The user with that email has already been registered.")
		default:
			log.Printf("signupPOST:CreateUser:%s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
			return
		}
		w.WriteHeader(http.StatusBadRequest) // 400
		form.ErrLog("signupPOST:")

		h.renderPage(w, "signup.html", &render.Data{
			Form: form,
		})
		return
	}

	http.Redirect(w, r, "/signin", http.StatusSeeOther) // 303
}
