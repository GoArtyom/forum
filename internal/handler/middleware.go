package handler

import (
	"context"
	"net/http"
	"time"

	"forum/pkg"
)

type conKay string

var keyUser = conKay("user")

func (h Handler) sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := pkg.GetCookie(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		session, err := h.service.GetSessionByUUID(cookie.Value)
		if err != nil {
			pkg.DeleteCookie(w)
			next.ServeHTTP(w, r)
			return
		}
		if session.ExpireAt.Before(time.Now()) {
			pkg.DeleteCookie(w)
			next.ServeHTTP(w, r)
			return
		}
		user, err := h.service.GetUserByUserId(session.User_id)
		if err != nil {
			pkg.DeleteCookie(w)
			h.service.DeleteSessionByUUID(cookie.Value)
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), keyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h Handler) authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := h.getUserFromContext(r)
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusSeeOther) // 303
			return
		}
		next.ServeHTTP(w, r)
	})
}
