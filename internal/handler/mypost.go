package handler

import (
	"forum/internal/models"
	"log"
	"net/http"
)

// GET
func (h Handler) myPostsGET(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/myposts" {
		log.Printf("myPostsGET: not found %s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}
	if r.Method != http.MethodGet {
		log.Printf("myPostsGET: method not allowed %s\n", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed) // 405
		return
	}
	user := h.getUserFromContext(r)
	posts, err := h.service.GetPostsByUserId(user.Id)
	//validate if a user has no posts
	if err != nil {
		log.Printf("myPostsGET: GetPostsByUserId: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}
	err = h.template.ExecuteTemplate(w, "index.html", models.Data{
		User:  user,
		Posts: posts,
	})

	if err != nil {
		log.Printf("myPostsGET: ExecuteTemplate %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
	}
}
