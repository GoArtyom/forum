package handler

import (
	"database/sql"
	"fmt"
	"forum/internal/models"
	"forum/internal/render"
	"forum/pkg/form"
	"log"
	"net/http"
)

// path /post/update?id=1"
func (h *Handler) updatePostGET_POST(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/update" {
		log.Printf("updatePostGET_POST:StatusNotFound:%s\n", r.URL.Path)
		h.renderError(w, http.StatusNotFound) // 404
		return
	}

	user := h.getUserFromContext(r)

	switch r.Method {
	case "GET":

		postId, err := h.getIntFromForm(r.URL.Query().Get("id"))
		if err != nil {
			log.Printf("updatePostGET:getIntFromForm():%s\n", err.Error())
			h.renderError(w, http.StatusBadRequest) // 400
			return
		}

		post, err := h.service.Post.GetById(postId)
		if err != nil {
			log.Printf("updatePostGET:GetById:%s\n", err.Error())
			if err == sql.ErrNoRows {
				h.renderError(w, http.StatusBadRequest) // 400
				return
			}
			h.renderError(w, http.StatusInternalServerError) // 500
			return
		}

		if post.UserId != user.Id {
			log.Printf("updatePostGET:Forbidden author:%d, user:%d", post.UserId, user.Id)
			h.renderError(w, http.StatusForbidden) // 403
			return
		}
		categories, err := h.service.Category.GetAll()
		if err != nil {
			log.Printf("updatePostGET:GetAll:%s\n", err.Error())
			h.renderError(w, http.StatusInternalServerError) // 500
			return
		}
		h.renderPage(w, "updatePost.html", &render.UpdatePost{
			User:       user,
			Post:       post,
			Categories: categories,
		})

	case "POST":
		if err := r.ParseMultipartForm(int64(21 << 20)); err != nil {
			log.Printf("updatePostPOST:ParseForm:%s\n", err.Error())
			h.renderError(w, http.StatusBadRequest) // 400
			return
		}
		form := form.New(r)
		_, handlerFile, err := r.FormFile("img")
		if err != nil && err != http.ErrMissingFile {
			log.Printf("updatePostPOST:FormFile:%s\n", err.Error())
			h.renderError(w, http.StatusInternalServerError) // 500
			return
		}

		if handlerFile != nil {
			form.ErrImg(handlerFile)
		}

		postId, err := h.getIntFromForm(r.Form.Get("id"))
		if err != nil {
			log.Printf("updatePostPOST:getIntFromForm():%s\n", err.Error())
			h.renderError(w, http.StatusBadRequest) // 400
			return
		}

		getCategories := r.PostForm["categories"]
		if len(getCategories) == 0 {
			form.Errors["categories"] = append(form.Errors["categories"], "You need to select at least one category.")
		}
		form.ErrEmpty("title", "content")
		form.ErrLengthMax("title", 50)
		form.ErrLengthMin("title", 5)
		form.ErrLengthMax("content", 5000)

		if len(form.Errors) != 0 {
			w.WriteHeader(http.StatusBadRequest) // 400
			form.ErrLog("updatePostPOST:")

			categories, err := h.service.Category.GetAll()
			if err != nil {
				log.Printf("updatePostPOST:GetAll:%s\n", err.Error())
				h.renderError(w, http.StatusInternalServerError) // 500
				return
			}

			h.renderPage(w, "updatePost.html", &render.CreatePost{
				User:       user,
				Categories: categories,
				Form:       form,
			})
			return
		}

		newPost := &models.UpdatePost{
			UserId:     user.Id,
			PostId:     postId,
			Title:      r.Form.Get("title"),
			Content:    r.Form.Get("content"),
			Categories: getCategories,
		}

		err = h.service.Post.UpdateById(newPost)
		if err != nil {
			log.Printf("updatePostPOST:UpdateById:%s\n", err.Error())
			if err == models.ErrPost {
				h.renderError(w, http.StatusBadRequest) // 400
				return
			}
			h.renderError(w, http.StatusInternalServerError) // 500
			return
		}

		if handlerFile != nil {
			newImage := &models.CreateImage{
				Header: handlerFile,
				PostId: postId,
			}

			err = h.service.Image.DeleteByPostId(postId)
			if err != nil && err != sql.ErrNoRows {
				log.Printf("updatePostPOST:DeleteByPostId: %s\n", err.Error())
				h.renderError(w, http.StatusInternalServerError) // 500
				return
			}

			err = h.service.Image.CreateByPostId(newImage)
			if err != nil {
				log.Printf("updatePostPOST:CreateByPostId: %s\n", err.Error())
				err = h.service.Post.DeleteById(&models.DeletePost{ServerErr: true})
				if err != nil {
					log.Printf("updatePostPOST:DeleteById: %s\n", err.Error())
				}
				h.renderError(w, http.StatusInternalServerError) // 500
				return
			}
		}

		http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postId), http.StatusSeeOther) // 303

	default:
		log.Printf("updatePostGET_POST:StatusMethodNotAllowed:%s\n", r.Method)
		h.renderError(w, http.StatusMethodNotAllowed) // 405
		return
	}
}
