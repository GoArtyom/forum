package handler

import (
	"log"
	"net/http"

	"forum/internal/render"
)

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		log.Printf("index:StatusNotFound:%s\n", r.URL.Path)
		h.renderError(w, http.StatusNotFound) // 404
		return
	}

	if r.Method != http.MethodGet {
		log.Printf("index:StatusMethodNotAllowed:%s\n", r.Method)
		h.renderError(w, http.StatusMethodNotAllowed) // 405
		return
	}

	user := h.getUserFromContext(r)

	posts, err := h.service.GetAllPost()
	if err != nil {
		log.Printf("index:GetAllPost:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	categories, err := h.service.GetAllCategory()
	if err != nil {
		log.Printf("index:GetAllCategory:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	h.renderPage(w, "home.html", &render.Data{
		User:       user,
		Posts:      posts,
		Categories: categories,
	})
}
