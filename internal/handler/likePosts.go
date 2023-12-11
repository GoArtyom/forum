package handler

import (
	"log"
	"net/http"

	"forum/internal/models"
)

func (h *Handler) likePostsGET(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/likeposts" {
		log.Printf("likePostsGET: not found %s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}
	if r.Method != http.MethodGet {
		log.Printf("likePostsGET: method not allowed %s\n", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed) // 405
		return
	}

	user := h.getUserFromContext(r)

	posts, err := h.service.GetPostsByLike(user.Id)
	if err != nil {
		log.Printf("likePostsGET: get all post %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}

	err = h.template.ExecuteTemplate(w, "index.html", models.Data{
		User:       user,
		Posts:      posts,
	})

	if err != nil {
		log.Printf("likePostsGET: ExecuteTemplate %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
	}
}
