package handler

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"forum/internal/models"
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
		session, err := h.service.Session.GetByUUID(cookie.Value)
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
		user, err := h.service.User.GetById(session.User_id)
		if err != nil {
			pkg.DeleteCookie(w, "UUID")
			h.service.DeleteByUUID(cookie.Value)
			next.ServeHTTP(w, r)
			return
		}

		count, err := h.service.Notification.GetCountByAuthorId(user.Id)
		if err != nil {
			log.Printf("sessionMiddleware:Notification.GetCountByAuthorId:%s\n", err.Error())
			h.renderError(w, http.StatusInternalServerError) // 500
			return
		}
		
		user.CountNotice = count

		ctx := context.WithValue(r.Context(), keyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) authUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := h.getUserFromContext(r)
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusSeeOther) // 303
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) authModerator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := h.getUserFromContext(r)
		if user.Role < models.ModeratorRole {
			h.renderError(w, http.StatusForbidden) // 403
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) authAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := h.getUserFromContext(r)
		if user.Role != models.AdminRole {
			h.renderError(w, http.StatusForbidden) // 403
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
				log.Println("recoverPanic:", err)
				w.Header().Set("Connection", "close")
				h.renderError(w, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
