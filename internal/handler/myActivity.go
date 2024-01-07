package handler

import (
	"forum/internal/render"
	"log"
	"net/http"
)

func (h *Handler) myActivityGET(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/myactivity" {
		log.Printf("myActivityGET:StatusNotFound:%s\n", r.URL.Path)
		h.renderError(w, http.StatusNotFound) // 404
		return
	}

	if r.Method != http.MethodGet {
		log.Printf("myActivityGET:NotMethodGet:%s\n", r.Method)
		h.renderError(w, http.StatusMethodNotAllowed) // 405
		return
	}

	user := h.getUserFromContext(r)

	myPosts, err := h.service.Post.GetAllByUserId(user.Id)
	if err != nil {
		log.Printf("myActivityGET:GetByUserId:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	myComments, err := h.service.Comment.GetAllByUserId(user.Id)
	if err != nil {
		log.Printf("myActivityGET:Comment.GetByUserId:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	votesPosts, err := h.service.Post.GetAllByUserVote(user.Id)
	if err != nil {
		log.Printf("myActivityGET:Post.GetByUserVote:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	votesComments, err := h.service.Comment.GetAllByUserVote(user.Id)
	if err != nil {
		log.Printf("myActivityGET:Comment.GetByUserVote:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	h.renderPage(w, "activity.html", &render.MyActivity{
		User:          user,
		MyPosts:       myPosts,
		MyComments:    myComments,
		VotesPosts:    votesPosts,
		VotesComments: votesComments,
	})

}
