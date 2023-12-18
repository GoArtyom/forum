package handler

import (
	"log"
	"net/http"

	"forum/internal/render"
)

// GET
func (h *Handler) myPostsGET(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/myposts" {
		log.Printf("myPostsGET:StatusNotFound:%s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}
	if r.Method != http.MethodGet {
		log.Printf("myPostsGET:StatusMethodNotAllowed:%s\n", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed) // 405
		return
	}
	user := h.getUserFromContext(r)
	posts, err := h.service.GetPostsByUserId(user.Id)
	// validate if a user has no posts
	if err != nil {
		log.Printf("myPostsGET:GetPostsByUserId:%s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}

	h.renderPage(w, "home.html", &render.Data{
		User:  user,
		Posts: posts,
	})
}
