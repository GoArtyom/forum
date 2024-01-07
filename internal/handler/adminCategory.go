package handler

import (
	"log"
	"net/http"

	"forum/internal/models"
	"forum/internal/render"
	"forum/pkg/form"
)

func (h *Handler) adminCategoriesDELETE(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/admin/categories/delete" {
		log.Printf("adminCategoriesDELETE:StatusNotFound:%s\n", r.URL.Path)
		h.renderError(w, http.StatusNotFound) // 404
		return
	}

	// delete
	if r.Method != http.MethodPost {
		log.Printf("amdinCategoriesDELETE:StatusMethodNotAllowed:%s\n", r.Method)
		h.renderError(w, http.StatusMethodNotAllowed) // 405
		return
	}

	category := r.PostFormValue("category")

	err := h.service.Category.DeleteByName(category)
	if err != nil {
		log.Printf("adminCategoriesDELETE:DeleteByName:%s\n", err.Error())
		if err == models.ErrCategory {
			h.renderError(w, http.StatusBadRequest) // 400
			return
		}
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther) // 303
}

/* endpoint: /admin/categories/create */
func (h *Handler) adminCategoriesCREATE(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/admin/categories/create" {
		log.Printf("adminCategoriesCREATE:StatusNotFound:%s\n", r.URL.Path)
		h.renderError(w, http.StatusNotFound) // 404
		return
	}

	if r.Method != http.MethodPost {
		log.Printf("adminCategoriesCREATE:StatusMethodNotAllowed:%s\n", r.Method)
		h.renderError(w, http.StatusMethodNotAllowed) // 405
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("adminCategoriesCREATE:ParseForm:%s\n", err.Error())
		h.renderError(w, http.StatusBadRequest) // 400
		return
	}

	form := form.New(r)
	form.ErrEmpty("category")

	if len(form.Errors) != 0 {
		w.WriteHeader(http.StatusBadRequest) // 400
		form.ErrLog("amdinCategoriesPOST:")
		admin := h.getUserFromContext(r)

		categories, err := h.service.Category.GetAll()
		if err != nil {
			log.Printf("adminGET:GetAllCategories:%s\n", err.Error())
			h.renderError(w, http.StatusInternalServerError) // 500
			return
		}

		reports, err := h.service.Report.GetAll()
		if err != nil {
			log.Printf("adminGET:GetAllReports:%s\n", err.Error())
			h.renderError(w, http.StatusInternalServerError) // 500
			return
		}

		users, err := h.service.User.GetAll()
		if err != nil {
			log.Printf("adminGET:GetAlUsersl:%s\n", err.Error())
			h.renderError(w, http.StatusInternalServerError) // 500
			return
		}

		h.renderPage(w, "admin.html", &render.AdminData{
			User:       admin,
			Reports:    reports,
			Users:      users,
			Categories: categories,
			Form:       form,
		})
		return
	}

	err := h.service.Category.Create(r.Form.Get("category"))
	if err != nil {
		// обработка ошибки если категория уже есть
		log.Printf("amdinCategoriesPOST:Create:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther) // 303
}
