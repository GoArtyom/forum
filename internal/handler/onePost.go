package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"forum/pkg/data"
)

// GET
func (h *Handler) onePostGET(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	if !strings.HasPrefix(r.URL.Path, "/post/") {
		log.Printf("onePostGET:StatusNotFound:%s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}
	if r.Method != http.MethodGet {
		log.Printf("onePostGET:StatusMethodNotAllowed:%s\n", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed) // 405
		return
	}
	postId, err := h.getPostIdFromURL(r.URL.Path)
	if err != nil {
		log.Printf("onePostGET:getPostIdFromURL:%s: %s\n", r.URL.Path, err.Error())
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}

	post, err := h.service.GetPostById(postId)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("onePostGET:GetPostById:post not found:%s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
			return
		}
		log.Printf("onePostGET:GetPostById:%s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}

	comments, err := h.service.GetAllCommentByPostId(post.PostId)
	if err != nil {
		log.Printf("onePostGET:GetAllCommentByPostId:%s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}
	user := h.getUserFromContext(r)

	err = h.template.ExecuteTemplate(w, "post.html", &data.Data{
		Post:     post,
		Comments: comments,
		User:     user,
	})

	if err != nil {
		log.Printf("onePostGET:ExecuteTemplate:%s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
	}
}
