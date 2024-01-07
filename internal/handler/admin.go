package handler

import (
	"log"
	"net/http"

	"forum/internal/models"
	"forum/internal/render"
)

func (h *Handler) adminGET(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/admin" {
		log.Printf("adminGET:StatusNotFound:%s\n", r.URL.Path)
		h.renderError(w, http.StatusNotFound) // 404
		return
	}
	if r.Method != http.MethodGet {
		log.Printf("adminGET:StatusMethodNotAllowed:%s\n", r.Method)
		h.renderError(w, http.StatusMethodNotAllowed) // 405
		return
	}

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
		User:             admin,
		Reports:          reports,
		Users:            users,
		RequestModerator: h.service.FilterByRole(users, models.ConsiderationModerator),
		Categories:       categories,
	})
}
