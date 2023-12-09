package handler

import (
	"forum/internal/models"
	"log"
	"net/http"
)

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		log.Printf("index: not found %s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}
	
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed) // 405
		return
	}

	user := h.getUserFromContext(r)

	posts, err := h.service.GetAllPost()
	if err != nil {
		log.Printf("index: get all post %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
	}

	categories, err := h.service.GetAllCategory()
	if err != nil {
		log.Printf("index: get all category %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
	}

	err = h.template.ExecuteTemplate(w, "index.html", models.Data{
		User:       user,
		Posts:      posts,
		Categories: categories,
	})

	if err != nil {
		log.Printf("index: ExecuteTemplate %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
	}
}
