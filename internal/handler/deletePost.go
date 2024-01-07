package handler

import (
	"log"
	"net/http"

	"forum/internal/models"
)

func (h *Handler) deletePostDELETE(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/delete" {
		log.Printf("deletePostDELETE:StatusNotFound:%s\n", r.URL.Path)
		h.renderError(w, http.StatusNotFound) // 404
		return
	}

	// DELETE
	if r.Method != http.MethodPost {
		log.Printf("deletePostDELETE:StatusMethodNotAllowed:%s\n", r.Method)
		h.renderError(w, http.StatusMethodNotAllowed) // 405
		return
	}

	postId, err := h.getIntFromForm(r.PostFormValue("id"))
	if err != nil {
		log.Printf("deletePostDELETE:getIntFromForm:%s: %s\n", r.URL.Path, err.Error())
		h.renderError(w, http.StatusNotFound) // 404
		return
	}

	user := h.getUserFromContext(r)

	formDel := &models.DeletePost{
		PostId:   postId,
		UserId:   user.Id,
		UserRole: user.Role,
	}

	err = h.service.Post.DeleteById(formDel)
	if err != nil {
		log.Printf("deletePostDELETE:DeleteById:%s\n", err.Error())
		if err == models.ErrPost {
			h.renderError(w, http.StatusBadRequest) // 400
			return
		}
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther) // 303
}
