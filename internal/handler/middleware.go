package handler

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"forum/pkg"
)

type conKay string

var keyUser = conKay("user")

func (h *Handler) sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := pkg.GetCookie(r, "UUID")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		session, err := h.service.GetSessionByUUID(cookie.Value)
		if err != nil {
			pkg.DeleteCookie(w, "UUID")
			next.ServeHTTP(w, r)
			return
		}
		if session.ExpireAt.Before(time.Now()) {
			pkg.DeleteCookie(w, "UUID")
			next.ServeHTTP(w, r)
			return
		}
		user, err := h.service.GetUserByUserId(session.User_id)
		if err != nil {
			pkg.DeleteCookie(w, "UUID")
			h.service.DeleteSessionByUUID(cookie.Value)
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), keyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := h.getUserFromContext(r)
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusSeeOther) // 303
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) limit(rate, cap float64, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Printf("limit:SplitHostPort: %s\n", err.Error())
			h.renderError(w, http.StatusInternalServerError)
			return
		}

		limiter := getVisitor(ip, rate, cap)
		if !limiter.Take(1) {
			addBlockList(ip)
			log.Printf("limit: TooManyRequests: %s\n", ip)
			h.renderError(w, http.StatusTooManyRequests) // 429
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				h.renderError(w, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
