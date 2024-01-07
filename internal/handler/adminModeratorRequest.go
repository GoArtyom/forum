package handler

import (
	"forum/internal/models"
	"log"
	"net/http"
)

func (h *Handler) adminModeratorRequestPATCH(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/admin/moderator-request" {
		log.Printf("adminModeratorRequestPATCH:StatusNotFound:%s\n", r.URL.Path)
		h.renderError(w, http.StatusNotFound) // 404
		return
	}
	if r.Method != http.MethodPost {
		log.Printf("adminModeratorRequestPATCH:StatusMethodNotAllowed:%s\n", r.Method)
		h.renderError(w, http.StatusMethodNotAllowed) // 405
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("adminModeratorRequestPATCH:ParseForm:%s\n", err.Error())
		h.renderError(w, http.StatusBadRequest) // 400
		return
	}

	userId, err := h.getIntFromForm(r.Form.Get("id"))
	if err != nil {
		log.Printf("adminModeratorRequestPATCH:getIntFromForm:%s: %s\n", r.URL.Path, err.Error())
		h.renderError(w, http.StatusBadRequest) // 404
		return
	}

	rele, err := h.getIntFromForm(r.Form.Get("role"))
	if err != nil {
		log.Printf("adminModeratorRequestPATCH:getIntFromForm:%s: %s\n", r.URL.Path, err.Error())
		h.renderError(w, http.StatusBadRequest) // 404
		return
	}

	user, err := h.service.User.GetById(userId)
	if err != nil {
		log.Printf("adminModeratorRequestPATCH:GetById:%s\n", err.Error())
		if err == models.ErrUser {
			h.renderError(w, http.StatusBadRequest) // 400
			return
		}
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	newRole := &models.UpdateRole{
		UserId:   userId,
		UserRole: user.Role,
		NewRole:  uint8(rele),
	}

	err = h.service.UpdateRoleById(newRole)
	if err != nil {
		log.Printf("moderatorRequestPATCH:UpdateUserRolsById:%s\n", err.Error())
		if err == models.ErrUpdateRole {
			h.renderError(w, http.StatusBadRequest) // 400
			return
		}

		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
