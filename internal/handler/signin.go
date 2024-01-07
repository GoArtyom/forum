package handler

import (
	"log"
	"net/http"

	"forum/internal/models"
	"forum/internal/render"
	"forum/pkg"
	"forum/pkg/form"
)

// GET
func (h *Handler) signinGET(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signin" {
		log.Printf("signinGET:StatusNotFound:%s\n", r.URL.Path)
		h.renderError(w, http.StatusNotFound) // 404
		return
	}
	if r.Method != http.MethodGet {
		log.Printf("signinGET:StatusMethodNotAllowed:%s\n", r.Method)
		h.renderError(w, http.StatusMethodNotAllowed) // 405
		return
	}

	h.renderPage(w, "signin.html", nil)
}

// POST
func (h *Handler) signinPOST(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/signin" {
		log.Printf("signinPOST:StatusNotFound:%s\n", r.URL.Path)
		h.renderError(w, http.StatusNotFound) // 404
		return
	}
	if r.Method != http.MethodPost {
		log.Printf("signinPOST:StatusMethodNotAllowed:%s\n", r.Method)
		h.renderError(w, http.StatusMethodNotAllowed) // 405
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Printf("signinPOST:ParseForm:%s\n", err.Error())
		h.renderError(w, http.StatusBadRequest) // 400
		return
	}

	// validate name/ email/ password
	form := form.New(r)

	form.ErrEmpty("email", "password")
	form.ErrLengthMin("email", 5)
	form.ErrLengthMax("email", 40)
	form.ErrLengthMin("password", 8)
	form.ErrLengthMax("password", 20)
	form.ValidEmail("email")
	form.ValidPassword("password")

	if len(form.Errors) != 0 {
		w.WriteHeader(http.StatusBadRequest) // 400
		form.ErrLog("signinPOST:")

		h.renderPage(w, "signin.html", &render.OnlyForm{
			Form: form,
		})
		return
	}

	user := &models.SignInUser{
		Email:    r.Form.Get("email"),
		Password: r.Form.Get("password"),
	}

	userId, err := h.service.User.SignIn(user)
	if err != nil {
		if err == models.ErrIncorData {
			w.WriteHeader(http.StatusBadRequest) // 400
			form.Errors["email"] = append(form.Errors["email"], "Email or password is incorrect.")
			form.ErrLog("signupPOST:")

			h.renderPage(w, "signin.html", &render.OnlyForm{
				Form: form,
			})
			return
		}
		log.Printf("signinPOST:SignInUser:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return

	}

	session, err := h.service.Session.Create(userId)
	if err != nil {
		log.Printf("signinPOST:Create:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	pkg.SetCookie(w, "UUID", session.UUID, session.ExpireAt)

	http.Redirect(w, r, "/", http.StatusSeeOther) // 303
}
