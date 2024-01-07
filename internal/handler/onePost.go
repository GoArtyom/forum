package handler

import (
	"database/sql"
	"log"
	"net/http"

	"forum/internal/render"
)

// GET
func (h *Handler) onePostGET(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post" {
		log.Printf("onePostGET:StatusNotFound:%s\n", r.URL.Path)
		h.renderError(w, http.StatusNotFound) // 404
		return
	}

	if r.Method != http.MethodGet {
		log.Printf("onePostGET:StatusMethodNotAllowed:%s\n", r.Method)
		h.renderError(w, http.StatusMethodNotAllowed) // 405
		return
	}

	id := r.URL.Query().Get("id")

	postId, err := h.getIntFromForm(id)
	if err != nil {
		log.Printf("onePostGET:getIntFromForm:%s: %s\n", r.URL.Path, err.Error())
		h.renderError(w, http.StatusNotFound) // 404
		return
	}

	post, err := h.service.Post.GetById(postId)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("onePostGET:GetById:post not found:%s\n", err.Error())
			h.renderError(w, http.StatusNotFound) // 400
			return
		}
		log.Printf("onePostGET:GetById:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	comments, err := h.service.Comment.GetAllByPostId(post.PostId)
	if err != nil {
		log.Printf("onePostGET:GetAllCommentByPostId:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	image, err := h.service.GetByPostId(post.PostId)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("onePostGET:GetByPostId:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}
	post.Image = image

	h.renderPage(w, "post.html", &render.OnePostData{
		Post:     post,
		Comments: comments,
		User:     h.getUserFromContext(r),
	})
}
