package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"forum/internal/models"
	"forum/internal/render"
	"forum/pkg/form"
)

func (h *Handler) createCommentPOST(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comment/create" {
		log.Printf("createCommentPOST:StatusNotFound:%s\n", r.URL.Path)
		h.renderError(w, http.StatusNotFound) // 404
		return
	}

	if r.Method != http.MethodPost {
		log.Printf("createCommentPOST:StatusMethodNotAllowed:%s\n", r.Method)
		h.renderError(w, http.StatusMethodNotAllowed) // 405
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("createCommentPOST:ParseForm:%s\n", err.Error())
		h.renderError(w, http.StatusBadRequest) // 400
		return
	}

	postId, err := h.getIntFromForm(r.Form.Get("post_id"))
	if err != nil {
		log.Printf("createCommentPOST:getIntFromForm():%s\n", err.Error())
		h.renderError(w, http.StatusBadRequest) // 400
		return
	}

	user := h.getUserFromContext(r)

	form := form.New(r)
	form.ErrEmpty("content")
	form.ErrLengthMin("content", 4)
	form.ErrLengthMax("content", 1000)

	if len(form.Errors) != 0 {
		form.ErrLog("createCommentPOST:")
		w.WriteHeader(http.StatusBadRequest)
		post, err := h.service.Post.GetById(postId)
		if err != nil {
			log.Printf("createCommentPOST:GetById:%s\n", err.Error())
			if err == sql.ErrNoRows {
				h.renderError(w, http.StatusBadRequest) // 400
				return
			}
			h.renderError(w, http.StatusInternalServerError) // 500
			return
		}

		comments, err := h.service.Comment.GetAllByPostId(post.PostId)
		if err != nil {
			log.Printf("createCommentPOST:GetAllCommentByPostId:%s\n", err.Error())
			h.renderError(w, http.StatusInternalServerError) // 500
			return
		}
		h.renderPage(w, "post.html", &render.OnePostData{
			User:     user,
			Post:     post,
			Comments: comments,
			Form:     form,
		})
		return
	}

	newComment := &models.CreateComment{
		PostId:   postId,
		Content:  r.Form.Get("content"),
		UserId:   user.Id,
		UserName: user.Name,
		CreateAt: time.Now(),
	}

	commentId, err := h.service.Comment.Create(newComment)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("createCommentPOST:CreateComment:post not found:%s\n", err.Error())
			h.renderError(w, http.StatusBadRequest) // 400
			return
		}
		log.Printf("createCommentPOST:CreateComment:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}
	newNoitification := &models.Notification{
		PostId:    postId,
		CommentId: commentId,
		UserId:    user.Id,
		UserName:  user.Name,
		Content:   r.Form.Get("content"),
		Type:      models.NoticeTypeComment,
		CreateAt:  time.Now(),
	}
	err = h.service.Notification.Create(newNoitification)
	if err != nil {
		log.Printf("createCommentPOST:CreateNotification:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postId), http.StatusSeeOther) // 303
}
