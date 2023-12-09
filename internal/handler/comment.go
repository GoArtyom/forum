package handler

import (
	"database/sql"
	"forum/internal/models"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (h Handler) createCommentPOST(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comment/create" {
		log.Printf("createCommentPOST: not found %s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}

	if r.Method != http.MethodPost {
		log.Printf("createCommentPOST: method not allowed %s\n", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed) // 405
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("createCommentPOST: parse form %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}

	postId, err := strconv.Atoi(r.Form.Get("postId"))
	//validation post id
	if err != nil {
		log.Printf("createCommentPOST: invalid post id: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
		return
	}

	//validation content
	user := h.getUserFromContext(r)
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
			log.Printf("createCommentPOST: post not found: %s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
			return
		}
		log.Printf("createCommentPOST: create comment: %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}

	err = h.template.ExecuteTemplate(w, "index.html", newComment)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
	}

}
