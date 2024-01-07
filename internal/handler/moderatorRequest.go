package handler

import (
	"log"
	"net/http"

	"forum/internal/models"
)

func (h *Handler) moderatorRequestPATCH(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/moderator/request" {
		log.Printf("moderatorRequestPATCH:StatusNotFound:%s\n", r.URL.Path)
		h.renderError(w, http.StatusNotFound) // 404
		return
	}

	if r.Method != http.MethodPost {
		log.Printf("moderatorRequestPATCH:StatusMethodNotAllowed:%s\n", r.Method)
		h.renderError(w, http.StatusMethodNotAllowed) // 405
		return
	}

	user := h.getUserFromContext(r)

	newRole := &models.UpdateRole{
		UserId:   user.Id,
		UserRole: user.Role,
		NewRole:  models.ConsiderationModerator,
	}

	err := h.service.UpdateRoleById(newRole)
	if err != nil {
		log.Printf("moderatorRequestPATCH:UpdateUserRolsById:%s\n", err.Error())

		if err == models.ErrUpdateRole {
			h.renderError(w, http.StatusBadRequest) // 400
			return
		}

		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther) // 303
}
