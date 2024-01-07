package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"forum/internal/models"
)

func (h *Handler) createCommentVotePOST(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comment/vote/create" {
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

	postId, err := h.getIntFromForm(r.Form.Get("post_id"))
	if err != nil {
		log.Printf("createCommentVotePOST:getIntFromForm():%s\n", err.Error())
		h.renderError(w, http.StatusBadRequest) // 400
		return
	}

	commentId, err := h.getIntFromForm(r.Form.Get("comment_id"))
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

	forNotic, err := h.service.CommentVote.Create(newVote)
	if err != nil {
		log.Printf("createCommentVotePOST:Create:%s\n", err.Error())
		if err.Error() == models.IncorRequest {
			h.renderError(w, http.StatusBadRequest) // 400
			return
		}
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	if forNotic%2 != 0 {
		//delete notic for vote
		formDel := &models.DeleteNotification{
			PostId: postId,
			UserId: user.Id,
			Type:   models.NoticeTypeCommentVote,
			Method: models.DelNoticByUser,
		}

		err = h.service.Notification.Delete(formDel)
		if err != nil {
			log.Printf("createPostVotePOST:Notification.Delete:%s\n", err.Error())
			if err == models.ErrNotification {
				h.renderError(w, http.StatusBadRequest) // 400
				return
			}
			h.renderError(w, http.StatusInternalServerError) // 500
			return
		}
	}

	if forNotic >= models.VoteSignalCreate {
		//create notic for vote
		newNotic := &models.Notification{
			PostId:    postId,
			CommentId: commentId,
			UserId:    user.Id,
			UserName:  user.Name,
			Vote:      vote,
			Type:      models.NoticeTypeCommentVote,
			CreateAt:  time.Now(),
		}

		err = h.service.Notification.Create(newNotic)
		if err != nil {
			log.Printf("createPostVotePOST:Notification.Create:%s\n", err.Error())
			h.renderError(w, http.StatusInternalServerError) // 500
			return
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postId), http.StatusSeeOther) // 303
}
