package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"forum/internal/models"
)

// GET
func (h Handler) onePostGET(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/post/") {
		log.Printf("onePostGET: not found %s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}
	if r.Method != http.MethodGet {
		log.Printf("onePostGET: method not allowed %s\n", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed) // 405
		return
	}
	postId, err := h.getPostIdFromURL(r.URL.Path)
	if err != nil {
		log.Printf("onePostGET: not found %s: %s\n", r.URL.Path, err.Error())
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}
	post, err := h.service.GetPostById(postId)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("onePostGET: post not found: %s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
			return
		}
		log.Printf("onePostGET: get post by id: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}
	
	err = h.template.ExecuteTemplate(w, "index.html", post)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
	}
}

// POST
func (h Handler) createPostPOST(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/create" {
		log.Printf("createPostPOST: not found %s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}
	if r.Method != http.MethodPost {
		log.Printf("createPostPOST: method not allowed %s\n", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed) // 405
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Printf("createPostPOST: parse form %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}
	// validate title/ content/
	categories := r.Form["categories"]
	if len(categories) == 0 {
		log.Println("createPostPOST: incorect len categories")
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
		log.Printf("createPostPOST: create post: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther) // 303
}
