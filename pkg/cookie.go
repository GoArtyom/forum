package pkg

import (
	"net/http"
	"time"
)

func SetCookie(w http.ResponseWriter, value string, expire_at time.Time) {
	cookie := &http.Cookie{
		Name:     "UUID",
		Value:    value,
		Path:     "/",
		Expires:  expire_at,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func GetCookie(r *http.Request) (*http.Cookie, error) {
	cookie, err := r.Cookie("UUID")
	if err != nil {
		return nil, err
	}
	return cookie, nil
}

func DeleteCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "UUID",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		MaxAge:   -1,
	}
	http.SetCookie(w, cookie)
}

func SetStateCookie(w http.ResponseWriter, state string) {
	cookie := &http.Cookie{
		Name:     "state",
		Value:    state,
		Path:     "/",
		MaxAge:   int(time.Hour.Seconds()),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}
