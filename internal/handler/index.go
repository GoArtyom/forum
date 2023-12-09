package handler

import (
	"fmt"
	"log"
	"net/http"
)

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		log.Printf("index: not found %s\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed) // 405
		return
	}
	// get post
	// get category
	// get userByContext
	err := h.template.ExecuteTemplate(w, "index.html", fmt.Sprintf("Path:%s\nMethod:%s", r.URL.Path, r.Method))
	if err != nil {
		log.Printf("index: ExecuteTemplate %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
	}
}
