package handler

import (
	"log"
	"net/http"

	"forum/internal/models"
	"forum/internal/render"
)

func (h *Handler) myNotificationsGET(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/mynotifications" {
		log.Printf("myNotificationsGET:StatusNotFound:%s\n", r.URL.Path)
		h.renderError(w, http.StatusNotFound) // 404
		return
	}

	user := h.getUserFromContext(r)

	switch r.Method {

	case "GET":

		notifications, err := h.service.Notification.GetAllByAuthorId(user.Id)
		if err != nil {
			log.Printf("myNotificationsGET:GetAllByUserId:%s\n", err.Error())
			h.renderError(w, http.StatusInternalServerError) // 500
			return
		}

		h.renderPage(w, "notifications.html", &render.MyNotifications{
			User:          user,
			Notifications: notifications,
		})

	case "POST":

		formDel := &models.DeleteNotification{
			UserId: user.Id,
			Method: models.DelNoticByAuthorAll,
		}
		h.service.Notification.Delete(formDel)

		http.Redirect(w, r, "/", http.StatusSeeOther) // 303
	default:
		log.Printf("myNotificationsGET:NotMethodGet:%s\n", r.Method)
		h.renderError(w, http.StatusMethodNotAllowed) // 405
	}
}
