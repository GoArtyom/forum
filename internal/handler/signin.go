package handler

import (
	"fmt"
	"log"
	"net/http"

	"forum/internal/models"
	"forum/pkg"
	"forum/pkg/data"
)

// GET
func (h *Handler) signinGET(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	if r.URL.Path != "/signin" {
		log.Printf("signinGET:StatusNotFound:%s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}
	if r.Method != http.MethodGet {
		log.Printf("signinGET:StatusMethodNotAllowed:%s\n", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed) // 405
		return
	}
	err := h.template.ExecuteTemplate(w, "signin.html", nil)
	if err != nil {
		log.Printf("signinGET:ExecuteTemplate:%s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
	}
}

// POST
func (h *Handler) signinPOST(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/signin" {
		log.Printf("signinPOST:StatusNotFound:%s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}
	if r.Method != http.MethodPost {
		log.Printf("signinPOST:StatusMethodNotAllowed:%s\n", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed) // 405
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Printf("signinPOST:ParseForm:%s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}

	// validate name/ email/ password
	data := new(data.Data)
	data.Errors = map[string][]string{}
	data.ErrEmpty(r, "email", "password")
	data.ErrLengthMin(r, "email", 5)
	data.ErrLengthMax(r, "email", 40)
	data.ErrLengthMin(r, "password", 8)
	data.ErrLengthMax(r, "password", 20)
	data.IsValid(r, "email", models.EmailRegexp)

	if len(data.Errors) != 0 {
		w.WriteHeader(http.StatusBadRequest) // 400
		data.ErrLog("signinPOST:")
		data.User = &models.User{
			Email:     r.Form.Get("email"),
			Password: r.Form.Get("password"),
		}
		err := h.template.ExecuteTemplate(w, "signin.html", data)
		if err != nil {
			log.Printf("signinPOST:ExecuteTemplate:%s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		}
		return
	}

	user := &models.SignInUser{
		Email:    r.Form.Get("email"),
		Password: r.Form.Get("password"),
	}
	userId, err := h.service.SignInUser(user)
	if err != nil {
		if err == models.ErrIncorData {
			w.WriteHeader(http.StatusBadRequest) // 400
			data.Errors["email"] = append(data.Errors["email"], "Email or password is incorrect.")
			data.ErrLog("signupPOST:")
			data.User = &models.User{
				Email:    r.Form.Get("email"),
				Password: r.Form.Get("password"),
			}

			err := h.template.ExecuteTemplate(w, "signin.html", data)
			if err != nil {
				log.Printf("signupPOST:ExecuteTemplate:%s\n", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
			}
			return
		}
		log.Printf("signinPOST:SignInUser:%s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return

	}

	session, err := h.service.CreateSession(userId)
	if err != nil {
		log.Printf("signinPOST:CreateSession:%s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}

	pkg.SetCookie(w, session.UUID, session.ExpireAt)

	http.Redirect(w, r, "/", http.StatusSeeOther) // 303
}
