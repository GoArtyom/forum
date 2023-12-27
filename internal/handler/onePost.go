package handler

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"forum/internal/render"
)

// GET
func (h *Handler) onePostGET(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/post/") {
		log.Printf("onePostGET:StatusNotFound:%s\n", r.URL.Path)
		h.renderError(w, http.StatusNotFound) // 404
		return
	}

	if r.Method != http.MethodGet {
		log.Printf("onePostGET:StatusMethodNotAllowed:%s\n", r.Method)
		h.renderError(w, http.StatusMethodNotAllowed) // 405
		return
	}

	postId, err := h.getPostIdFromURL(r.URL.Path)
	if err != nil {
		log.Printf("onePostGET:getPostIdFromURL:%s: %s\n", r.URL.Path, err.Error())
		h.renderError(w, http.StatusNotFound) // 404
		return
	}

	post, err := h.service.GetPostById(postId)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("onePostGET:GetPostById:post not found:%s\n", err.Error())
			h.renderError(w, http.StatusNotFound) // 400
			return
		}
		log.Printf("onePostGET:GetPostById:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	comments, err := h.service.GetAllCommentByPostId(post.PostId)
	if err != nil {
		log.Printf("onePostGET:GetAllCommentByPostId:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	image, err := h.service.GetImageByPostId(post.PostId)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("onePostGET:GetImageByPostId:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}
	post.Image = image

	h.renderPage(w, "post.html", &render.Data{
		Post:     post,
		Comments: comments,
		User:     h.getUserFromContext(r),
	})
}
