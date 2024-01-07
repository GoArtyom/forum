package handler

import (
	"fmt"
	"forum/internal/models"
	"forum/internal/render"
	"forum/pkg/form"
	"log"
	"net/http"
)

func (h *Handler) updateCommentGET_POST(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comment/update" {
		log.Printf("updateCommentGET_POST:StatusNotFound:%s\n", r.URL.Path)
		h.renderError(w, http.StatusNotFound) // 404
		return
	}

	user := h.getUserFromContext(r)

	switch r.Method {
	case "GET":

		commentId, err := h.getIntFromForm(r.URL.Query().Get("id"))
		if err != nil {
			log.Printf("updatePostGET:getIntFromForm():%s\n", err.Error())
			h.renderError(w, http.StatusBadRequest) // 400
			return
		}

		comment, err := h.service.Comment.GetById(commentId)
		if err != nil {
			log.Printf("updateCommentGET:GetById:%s\n", err.Error())
			if err == models.ErrComment {
				h.renderError(w, http.StatusBadRequest) // 400
				return
			}
			h.renderError(w, http.StatusInternalServerError) // 500
			return
		}
		if comment.UserId != user.Id {
			log.Printf("updateCommentGET:Forbidden author:%d, user:%d", comment.UserId, user.Id)
			h.renderError(w, http.StatusForbidden) // 403
			return
		}

		h.renderPage(w, "updateComment.html", &render.UpdateComment{
			User:    user,
			Comment: comment,
		})

	case "POST":
		if err := r.ParseForm(); err != nil {
			log.Printf("updateCommentPOST:ParseForm:%s\n", err.Error())
			h.renderError(w, http.StatusBadRequest) // 400
			return
		}
		commentId, err := h.getIntFromForm(r.Form.Get("id"))
		if err != nil {
			log.Printf("updateCommentPOST:getIntFromForm:%s\n", err.Error())
			h.renderError(w, http.StatusBadRequest) // 400
			return
		}

		postId, err := h.getIntFromForm(r.Form.Get("post_id"))
		if err != nil {
			log.Printf("updateCommentPOST:getIntFromForm:%s\n", err.Error())
			h.renderError(w, http.StatusBadRequest) // 400
			return
		}

		form := form.New(r)
		form.ErrEmpty("content")
		form.ErrLengthMin("content", 4)
		form.ErrLengthMax("content", 1000)

		if len(form.Errors) != 0 {
			form.ErrLog("updateCommentPOST:")
			w.WriteHeader(http.StatusBadRequest)

			comment, err := h.service.Comment.GetById(commentId)
			if err != nil {
				log.Printf("updateCommentPOST:GetById:%s\n", err.Error())
				if err == models.ErrComment {
					h.renderError(w, http.StatusBadRequest) // 400
					return
				}
				h.renderError(w, http.StatusInternalServerError) // 500
				return
			}
			h.renderPage(w, "updateComment.html", &render.UpdateComment{
				User:    user,
				Comment: comment,
				Form:    form,
			})
			return
		}

		updateComment := &models.UpdateComment{
			Id:      commentId,
			UserId:  user.Id,
			Content: r.Form.Get("content"),
		}

		err = h.service.Comment.UpdateById(updateComment)
		if err != nil {
			log.Printf("updateCommentPOST:UpdateById:%s\n", err.Error())
			if err == models.ErrComment {
				h.renderError(w, http.StatusBadRequest) // 400
				return
			}
			h.renderError(w, http.StatusInternalServerError) // 500
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postId), http.StatusSeeOther) // 303

	default:
		log.Printf("updateCommentGET_POST:StatusMethodNotAllowed:%s\n", r.Method)
		h.renderError(w, http.StatusMethodNotAllowed) // 405
		return
	}
}
