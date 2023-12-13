package handler

import (
	"fmt"
	"log"
	"net/http"

	"forum/internal/models"
	"forum/pkg/data"
)

// GET
func (h *Handler) signupGET(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
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
	err := h.template.ExecuteTemplate(w, "signup.html", nil)
	if err != nil {
		log.Printf("signupGET:ExecuteTemplate:%s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
	}
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
	data := new(data.Data)
	data.Errors = map[string][]string{}
	data.ErrEmpty(r, "name", "email", "password")
	data.ErrLengthMin(r, "name", 5)
	data.ErrLengthMax(r, "name", 20)
	data.ErrLengthMin(r, "email", 5)
	data.ErrLengthMax(r, "email", 40)
	data.ErrLengthMin(r, "password", 8)
	data.ErrLengthMax(r, "password", 20)
	data.IsValid(r, "email", models.EmailRegexp)
	// data.IsValid(r, "password", models.PasswordRegexp)

	if len(data.Errors) != 0 {
		w.WriteHeader(http.StatusBadRequest) // 400
		data.ErrLog("signupPOST:")

		data.User = &models.User{
			Name:     r.Form.Get("name"),
			Email:    r.Form.Get("email"),
			Password: r.Form.Get("password"),
		}

		err := h.template.ExecuteTemplate(w, "signup.html", data)
		if err != nil {
			log.Printf("signupPOST:ExecuteTemplate:%s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		}
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
			w.WriteHeader(http.StatusBadRequest) // 400
			data.Errors["name"] = append(data.Errors["name"], "The user with that name has already been registered.")
			data.ErrLog("signupPOST:")
			data.User = &models.User{
				Name:     r.Form.Get("name"),
				Email:    r.Form.Get("email"),
				Password: r.Form.Get("password"),
			}

			err := h.template.ExecuteTemplate(w, "signup.html", data)
			if err != nil {
				log.Printf("signupPOST:ExecuteTemplate:%s\n", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
			}
			return
		case models.UniqueEmail:
			w.WriteHeader(http.StatusBadRequest) // 400
			data.Errors["email"] = append(data.Errors["email"], "The user with that email has already been registered.")
			data.ErrLog("signupPOST:")
			data.User = &models.User{
				Name:     r.Form.Get("name"),
				Email:    r.Form.Get("email"),
				Password: r.Form.Get("password"),
			}

			err := h.template.ExecuteTemplate(w, "signup.html", data)
			if err != nil {
				log.Printf("signupPOST:ExecuteTemplate:%s\n", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
			}
			return
		default:
			log.Printf("signupPOST:CreateUser:%s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
			return
		}
	}
	http.Redirect(w, r, "/signin", http.StatusSeeOther) // 303
}
