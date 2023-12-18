package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"forum/internal/models"
	"forum/internal/render"
	"forum/pkg/form"
)

func (h *Handler) createCommentPOST(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comment/create" {
		log.Printf("createCommentPOST:StatusNotFound:%s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}

	if r.Method != http.MethodPost {
		log.Printf("createCommentPOST:StatusMethodNotAllowed:%s\n", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed) // 405
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("createCommentPOST:ParseForm:%s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}

	postId, err := h.getIntFromForm(r, "post_id")
	if err != nil {
		log.Printf("createCommentPOST:getIntFromForm():%s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
		return
	}

	user := h.getUserFromContext(r)

	form := form.New(r)
	form.ErrEmpty("content")
	form.ErrLengthMin("content", 5)
	form.ErrLengthMax("content", 1000)

	if len(form.Errors) != 0 {
		form.ErrLog("createCommentPOST:")
		w.WriteHeader(http.StatusBadRequest)
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
		h.renderPage(w, "post.html", &render.Data{
			User:     user,
			Post:     post,
			Comments: comments,
			Form:     form,
		})
		return
	}

	newComment := &models.CreateComment{
		PostId:   postId,
		Content:  r.Form.Get("content"),
		UserId:   user.Id,
		UserName: user.Name,
		CreateAt: time.Now(),
	}

	err = h.service.CreateComment(newComment)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("createCommentPOST:CreateComment:post not found:%s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
			return
		}
		log.Printf("createCommentPOST:CreateComment:%s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", postId), http.StatusSeeOther) // 303
}
