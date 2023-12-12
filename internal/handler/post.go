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

	// add like in post
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

	// add like in comment
	comments, err := h.service.GetAllCommentByPostId(post.PostId)
	if err != nil {
		log.Printf("onePostGET:GetAllCommentByPostId:%s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}
	err = h.template.ExecuteTemplate(w, "post.html", &models.Data{
		Post:     post,
		Comments: comments,
	})

	if err != nil {
		log.Printf("onePostGET:ExecuteTemplate:%s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
	}
}





func (h *Handler) createPostGET_POST(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	if r.URL.Path != "/post/create" {
		log.Printf("createPostGET_POST:StatusNotFound:%s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}
	switch r.Method {

	// POST
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			log.Printf("createPostPOST:ParseForm:%s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
			return
		}
		// validate title/ content/

		categories := r.Form["categories"]
		if len(categories) == 0 {
			log.Println("createPostPOST:incorect len categories")
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
			Categories: categories,
			CreateAt:   time.Now(),
		}
		id, err := h.service.CreatePost(newPost)
		if err != nil {
			log.Printf("createPostPOST:CreatePost:%s\n", err.Error())
			if err.Error() == models.IncorRequest {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
				return
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther) // 303

	// GET
	case http.MethodGet:
		user := h.getUserFromContext(r)

		categories, err := h.service.GetAllCategory()
		if err != nil {
			log.Printf("createPostGET:GetAllCategory:%s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
			return
		}

		err = h.template.ExecuteTemplate(w, "create.html", models.Data{
			User:       user,
			Categories: categories,
		})

		if err != nil {
			log.Printf("createPostGET:ExecuteTemplate:%s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		}

	default:
		log.Printf("createPostGET_POST:StatusMethodNotAllowed:%s\n", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed) // 405
		return
	}
}
