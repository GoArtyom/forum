package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"forum/internal/models"
)

func (h *Handler) createCommentVotePOST(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/comment/vote/create") {
		log.Printf("createCommentVotePOST:StatusNotFound:%s\n", r.URL.Path)
		h.renderError(w, http.StatusNotFound) // 404
		return
	}

	if r.Method != http.MethodPost {
		log.Printf("createCommentVotePOST:StatusMethodNotAllowed:%s\n", r.Method)
		h.renderError(w, http.StatusMethodNotAllowed) // 405
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("createCommentVotePOST:ParseForm:%s\n", err.Error())
		h.renderError(w, http.StatusBadRequest) // 400
		return
	}

	vote, err := h.getVote(r.Form.Get("vote"))
	if err != nil {
		log.Printf("createCommentVotePOST:getVote:%s\n", err.Error())
		h.renderError(w, http.StatusBadRequest) // 400
		return
	}

	postId, err := h.getIntFromForm(r, "post_id")
	if err != nil {
		log.Printf("createCommentVotePOST:getIntFromForm():%s\n", err.Error())
		h.renderError(w, http.StatusBadRequest) // 400
		return
	}

	commentId, err := h.getIntFromForm(r, "comment_id")
	if err != nil {
		log.Printf("createCommentVotePOST:getIntFromForm():%s\n", err.Error())
		h.renderError(w, http.StatusBadRequest) // 400
		return
	}

	user := h.getUserFromContext(r)

	newVote := &models.CommentVote{
		CommentId: commentId,
		UserId:    user.Id,
		Vote:      vote,
	}

	err = h.service.CommentVote.CreateCommentVote(newVote)
	if err != nil {
		log.Printf("createCommentVotePOST:CreatePostVote:%s\n", err.Error())
		if err.Error() == models.IncorRequest {
			h.renderError(w, http.StatusBadRequest) // 400
			return
		}
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", postId), http.StatusSeeOther) // 303
}
