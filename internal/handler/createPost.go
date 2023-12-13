package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"forum/internal/models"
	"forum/pkg/data"
)

func (h *Handler) createPostGET_POST(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	if r.URL.Path != "/post/create" {
		log.Printf("createPostGET_POST:StatusNotFound:%s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}

	user := h.getUserFromContext(r)

	switch r.Method {

	// POST
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			log.Printf("createPostPOST:ParseForm:%s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
			return
		}

		data := new(data.Data)
		data.Errors = map[string][]string{}
		getCategories := r.PostForm["categories"]
		if len(getCategories) == 0 {
			data.Errors["categories"] = append(data.Errors["categories"], "You need to select at least one category.")
		}
		data.ErrEmpty(r, "title", "content")
		data.ErrLengthMax(r, "title", 50)
		data.ErrLengthMin(r, "title", 10)
		data.ErrLengthMax(r, "content", 2000)

		if len(data.Errors) != 0 {
			w.WriteHeader(http.StatusBadRequest) // 400
			data.ErrLog("createPostPOST:")
		
			categories, err := h.service.GetAllCategory()
			if err != nil {
				log.Printf("createPostGET:GetAllCategory:%s\n", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
				return
			}

			data.User = user
			data.Categories = categories
			data.Post = &models.Post{
				Title:      r.Form.Get("title"),
				Content:    r.Form.Get("content"),
				Categories: getCategories,
			}

			err = h.template.ExecuteTemplate(w, "create.html", data)

			if err != nil {
				log.Printf("createPostGET:ExecuteTemplate:%s\n", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
			}
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
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
				return
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther) // 303

	// GET
	case http.MethodGet:

		categories, err := h.service.GetAllCategory()
		if err != nil {
			log.Printf("createPostGET:GetAllCategory:%s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
			return
		}

		err = h.template.ExecuteTemplate(w, "create.html", data.Data{
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
