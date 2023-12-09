package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"forum/internal/models"
)

func (h Handler) createPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/create" {
		log.Printf("createPost: not found %s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}
	if r.Method != http.MethodPost {
		log.Printf("createPost: method not allowed %s\n", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed) // 405
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Printf("createPost: parse form %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
	}
	// validate title/ content/
	categories := r.Form["categories"]
	if len(categories) == 0 {
		log.Println("createPost: incorect len categories")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
		return
	}
	// validate categories
	user := h.getUserFromContext(r)
	newPost := &models.CreatePost{
		Title:      r.Form.Get("title"),
		Content:    r.Form.Get("content"),
		UserId:     user.Id,
		UserName:   user.Name,
		Categories: &categories,
		CreateAt:   time.Now(),
	}
	id, err := h.service.CreatePost(newPost)
	if err != nil {
		log.Printf("createPost: create post: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther) // 303
}
