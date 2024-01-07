package handler

import (
	"forum/internal/models"
	"log"
	"net/http"
)

func (h *Handler) adminReportDELETE(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/admin/report" {
		log.Printf("adminReportDELETE:StatusNotFound:%s\n", r.URL.Path)
		h.renderError(w, http.StatusNotFound) // 404
		return
	}

	//delete
	if r.Method != http.MethodPost {
		log.Printf("adminReportDELETE:StatusMethodNotAllowed:%s\n", r.Method)
		h.renderError(w, http.StatusMethodNotAllowed) // 405
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("adminReportDELETE:ParseForm:%s\n", err.Error())
		h.renderError(w, http.StatusBadRequest) // 400
		return
	}

	reportId, err := h.getIntFromForm(r.Form.Get("id"))
	if err != nil {
		log.Printf("adminReportDELETE:getIntFromForm:%s: %s\n", r.URL.Path, err.Error())
		h.renderError(w, http.StatusBadRequest) // 404
		return
	}

	resp, err := h.getIntFromForm(r.Form.Get("resp"))
	if err != nil {
		log.Printf("adminReportDELETE:getIntFromForm:%s: %s\n", r.URL.Path, err.Error())
		h.renderError(w, http.StatusBadRequest) // 404
		return
	}

	err = h.service.Report.DeleteById(reportId, resp)
	if err != nil {
		log.Printf("adminReportDELETE:DeleteById:%s\n", err.Error())
		if err == models.ErrReport {
			h.renderError(w, http.StatusBadRequest) // 400
			return
		}
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
