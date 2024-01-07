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

/* endpoint = /post/reporting?id=3 */
func (h *Handler) reportingPostPOST(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/reporting" {
		log.Printf("reportingPostPOST:StatusNotFound:%s\n", r.URL.Path)
		h.renderError(w, http.StatusNotFound) // 404
		return
	}
	if r.Method != http.MethodPost {
		log.Printf("reportingPostPOST:StatusMethodNotAllowed:%s\n", r.Method)
		h.renderError(w, http.StatusMethodNotAllowed) // 405
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("reportingPostPOST:ParseForm:%s\n", err.Error())
		h.renderError(w, http.StatusBadRequest) // 400
		return
	}

	postId, err := h.getIntFromForm(r.Form.Get("id"))
	if err != nil {
		log.Printf("reportingPostPOST:getIntFromForm:%s: %s\n", r.URL.Path, err.Error())
		h.renderError(w, http.StatusNotFound) // 404
		return
	}

	user := h.getUserFromContext(r)

	form := form.New(r)
	form.ErrEmpty("report")
	form.ErrLengthMin("report", 4)
	form.ErrLengthMax("report", 100)

	if len(form.Errors) != 0 {

		form.ErrLog("reportingPostPOST:")

		w.WriteHeader(http.StatusBadRequest)

		post, err := h.service.Post.GetById(postId)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("reportingPostPOST:GetById:post not found:%s\n", err.Error())
				h.renderError(w, http.StatusBadRequest) // 400
				return
			}
			log.Printf("reportingPostPOST:GetById:%s\n", err.Error())
			h.renderError(w, http.StatusInternalServerError) // 500
			return
		}

		comments, err := h.service.Comment.GetAllByPostId(post.PostId)
		if err != nil {
			log.Printf("reportingPostPOST:GetAllCommentByPostId:%s\n", err.Error())
			h.renderError(w, http.StatusInternalServerError) // 500
			return
		}
		h.renderPage(w, "post.html", &render.OnePostData{
			User:     user,
			Post:     post,
			Comments: comments,
			Form:     form,
		})
		return
	}

	newReport := &models.CreateReport{
		PostId:        postId,
		Content:       r.Form.Get("report"),
		ModeratorId:   user.Id,
		ModeratorName: user.Name,
		CreateAt:      time.Now(),
	}
	err = h.service.Report.Create(newReport)

	if err != nil {
		log.Printf("reportingPostPOST:CreateReport:%s\n", err.Error())
		if err == models.ErrPost {
			h.renderError(w, http.StatusBadRequest) // 400
			return
		}
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postId), http.StatusSeeOther) // 303
}
