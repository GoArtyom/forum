package handler

import (
	"fmt"
	"log"
	"net/http"

	"forum/internal/models"
)

func (h *Handler) likePostsGET(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	if r.URL.Path != "/likeposts" {
		log.Printf("likePostsGET:StatusNotFound:%s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}
	if r.Method != http.MethodGet {
		log.Printf("likePostsGET:StatusMethodNotAllowed:%s\n", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed) // 405
		return
	}

	user := h.getUserFromContext(r)

	posts, err := h.service.GetPostsByLike(user.Id)
	if err != nil {
		log.Printf("likePostsGET:GetPostsByLike:%s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}

	err = h.template.ExecuteTemplate(w, "home.html", models.Data{
		User:  user,
		Posts: posts,
	})

	if err != nil {
		log.Printf("likePostsGET:ExecuteTemplate:%s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
	}
}
