package handler

import (
	"forum/internal/models"
	"log"
	"net/http"
)

// POST
func (h *Handler) filterPostsPOST(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/filterposts" {
		log.Printf("filterPostsPOST: not found %s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}

	if r.Method != http.MethodPost {
		log.Printf("filterPostsPOST: method not allowed %s\n", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed) // 405
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("filterPostsPOST: parse form %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}
	// validate title/ content/
	if len(r.Form["categories"]) != 1 {
		log.Printf("filterPostsPOST:bad request:len(categories) = %d\n", len(r.Form["categories"]))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
		return
	}

	posts, err := h.service.GetPostsByCategory(r.Form.Get("categories"))
	if err != nil {
		log.Printf("filterPostsPOST: GetPostsByCategory %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}
	categories, err := h.service.GetAllCategory()
	if err != nil {
		log.Printf("filterPostsPOST: GetPostsByCategory %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}

	user := h.getUserFromContext(r)

	err = h.template.ExecuteTemplate(w, "index.html", &models.Data{
		User:       user,
		Posts:      posts,
		Categories: categories,
	})

	if err != nil {
		log.Printf("index: ExecuteTemplate %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
	}
}
