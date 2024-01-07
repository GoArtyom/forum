package handler

import (
	"log"
	"net/http"

	"forum/internal/render"
)

func (h *Handler) indexGET(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		log.Printf("indexGET:StatusNotFound:%s\n", r.URL.Path)
		h.renderError(w, http.StatusNotFound) // 404
		return
	}

	if r.Method != http.MethodGet {
		log.Printf("indexGET:StatusMethodNotAllowed:%s\n", r.Method)
		h.renderError(w, http.StatusMethodNotAllowed) // 405
		return
	}

	user := h.getUserFromContext(r)

	posts, err := h.service.Post.GetAll()
	if err != nil {
		log.Printf("indexGET:GetAll:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	categories, err := h.service.Category.GetAll()
	if err != nil {
		log.Printf("indexGET:GetAll:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	h.renderPage(w, "home.html", &render.MainData{
		User:       user,
		Posts:      posts,
		Categories: categories,
	})
}
