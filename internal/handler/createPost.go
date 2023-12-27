package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"forum/internal/models"
	"forum/internal/render"
	"forum/pkg/form"
)

func (h *Handler) createPostGET_POST(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/create" {
		log.Printf("createPostGET_POST:StatusNotFound:%s\n", r.URL.Path)
		h.renderError(w, http.StatusNotFound) // 404
		return
	}

	user := h.getUserFromContext(r)

	switch r.Method {

	// POST
	case http.MethodPost:
		if err := r.ParseMultipartForm(int64(21 << 20)); err != nil {
			log.Printf("createPostPOST:ParseForm:%s\n", err.Error())
			h.renderError(w, http.StatusBadRequest) // 400
			return
		}
		form := form.New(r)
		_, handlerFile, err := r.FormFile("img")
		if err != nil && err != http.ErrMissingFile {
			log.Printf("createPostPOST:FormFile:%s\n", err.Error())
			h.renderError(w, http.StatusInternalServerError) // 500
			return
		}

		if handlerFile != nil {
			form.ErrImg(handlerFile)
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
			form.ErrLog("createPostPOST:")

			categories, err := h.service.GetAllCategory()
			if err != nil {
				log.Printf("createPostGET:GetAllCategory:%s\n", err.Error())
				h.renderError(w, http.StatusInternalServerError) // 500
				return
			}

			h.renderPage(w, "create.html", &render.Data{
				User:       user,
				Categories: categories,
				Form:       form,
			})
			return
		}

		newPost := &models.CreatePost{
			Title:      r.Form.Get("title"),
			Content:    r.Form.Get("content"),
			UserId:     user.Id,
			UserName:   user.Name,
			Categories: getCategories,
			CreateAt:   time.Now(),
		}

		id, err := h.service.CreatePost(newPost)
		if err != nil {
			log.Printf("createPostPOST:CreatePost:%s\n", err.Error())
			if err.Error() == models.IncorRequest {
				h.renderError(w, http.StatusBadRequest) // 400
				return
			}
			h.renderError(w, http.StatusInternalServerError) // 500
			return
		}

		newImage := &models.CreateImage{
			Header: handlerFile,
			PostId: id,
		}

		if handlerFile != nil {
			err = h.service.CreateImageByPostIt(newImage)
			if err != nil {
				log.Printf("createPostPOST:CreateImageByPostIt: %s\n", err.Error())
				err = h.service.DeletePostById(id)
				if err != nil {
					log.Printf("createPostPOST:DeletePostById: %s\n", err.Error())
				}
				h.renderError(w, http.StatusInternalServerError) // 500
				return
			}
		}

		http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther) // 303

	// GET
	case http.MethodGet:

		categories, err := h.service.GetAllCategory()
		if err != nil {
			log.Printf("createPostGET:GetAllCategory:%s\n", err.Error())
			h.renderError(w, http.StatusInternalServerError) // 500
			return
		}
		h.renderPage(w, "create.html", &render.Data{
			User:       user,
			Categories: categories,
		})

	default:
		log.Printf("createPostGET_POST:StatusMethodNotAllowed:%s\n", r.Method)
		h.renderError(w, http.StatusMethodNotAllowed) // 405
		return
	}
}
