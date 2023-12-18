package handler

import (
	"database/sql"
	"log"
	"net/http"

	"forum/internal/render"
)

// GET
func (h *Handler) filterPostsGET(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/filterposts" {
		log.Printf("filterPostsGET:StatusNotFound:%s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}

	if r.Method != http.MethodGet {
		log.Printf("filterPostsGET:StatusMethodNotAllowed:%s\n", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed) // 405
		return
	}
	category := r.URL.Query().Get("category")

	posts, err := h.service.GetPostsByCategory(category)
	if err != nil {
		log.Printf("filterPostsGET:GetPostsByCategory:%s\n", err.Error())
		if err == sql.ErrNoRows {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}
	categories, err := h.service.GetAllCategory()
	if err != nil {
		log.Printf("filterPostsGET:GetAllCategory:%s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}

	user := h.getUserFromContext(r)

	h.renderPage(w, "home.html", &render.Data{
		User:       user,
		Posts:      posts,
		Categories: categories,
	})
}
