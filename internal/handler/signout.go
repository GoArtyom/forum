package handler

import (
	"forum/pkg"
	"log"
	"net/http"
)

func (h *Handler) signoutPOST(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/signout" {
		log.Printf("signoutPOST: not found %s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}
	if r.Method != http.MethodPost {
		log.Printf("signoutPOST: method not allowed %s\n", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed) // 405
		return
	}
	
	session, err := pkg.GetCookie(r)
	if err != nil {
		log.Printf("signoutPOST:get cookie %s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}
	err = h.service.DeleteSessionByUUID(session.Value)
	if err != nil {
		log.Printf("signoutPOST:delete session by UUID %s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}
	pkg.DeleteCookie(w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
