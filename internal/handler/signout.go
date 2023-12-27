package handler

import (
	"log"
	"net/http"

	"forum/pkg"
)

func (h *Handler) signoutPOST(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/signout" {
		log.Printf("signoutPOST:StatusNotFound:%s\n", r.URL.Path)
		h.renderError(w, http.StatusNotFound) // 404
		return
	}
	if r.Method != http.MethodPost {
		log.Printf("signoutPOST:StatusMethodNotAllowed:%s\n", r.Method)
		h.renderError(w, http.StatusMethodNotAllowed) // 405
		return
	}

	session, err := pkg.GetCookie(r, "UUID")
	if err != nil {
		log.Printf("signoutPOST:GetCookie:%s\n", r.URL.Path)
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}
	err = h.service.DeleteSessionByUUID(session.Value)
	if err != nil {
		log.Printf("signoutPOST:DeleteSessionByUUID:%s\n", r.URL.Path)
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	pkg.DeleteCookie(w, "UUID")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
