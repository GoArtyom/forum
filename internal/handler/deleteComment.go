package handler

import (
	"fmt"
	"log"
	"net/http"

	"forum/internal/models"
)

func (h *Handler) deleteCommentDELETE(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comment/delete" {
		log.Printf("deleteCommentDELETE:StatusNotFound:%s\n", r.URL.Path)
		h.renderError(w, http.StatusNotFound) // 404
		return
	}

	//delete
	if r.Method != http.MethodPost {
		log.Printf("deleteCommentDELETE:StatusMethodNotAllowed:%s\n", r.Method)
		h.renderError(w, http.StatusMethodNotAllowed) // 405
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("deleteCommentDELETE:ParseForm:%s\n", err.Error())
		h.renderError(w, http.StatusBadRequest) // 404
		return
	}

	commentId, err := h.getIntFromForm(r.Form.Get("id"))
	if err != nil {
		log.Printf("deleteCommentDELETE:getIntFromForm:%s: %s\n", r.URL.Path, err.Error())
		h.renderError(w, http.StatusBadRequest) // 404
		return
	}

	postId, err := h.getIntFromForm(r.Form.Get("post_id"))
	if err != nil {
		log.Printf("deleteCommentDELETE:getIntFromForm:%s: %s\n", r.URL.Path, err.Error())
		h.renderError(w, http.StatusBadRequest) // 404
		return
	}
	/// проверка post_id

	user := h.getUserFromContext(r)

	formDel := &models.DeleteComment{
		CommentId: commentId,
		UserId:    user.Id,
		UserRole:  user.Role,
	}
	err = h.service.Comment.DeleteById(formDel)
	if err != nil {
		log.Printf("deleteCommentDELETE:DeleteById:%s\n", err.Error())
		if err == models.ErrComment {
			h.renderError(w, http.StatusBadRequest) // 400
			return
		}
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postId), http.StatusSeeOther) // 303
}
