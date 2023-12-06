package handler

import "net/http"

func (h Handler) InitRouters() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", h.index)

	mux.HandleFunc("/signin", h.signin)
	mux.HandleFunc("/signup", h.signup)

	mux.HandleFunc("/auth/signin", h.signinPost)
	mux.HandleFunc("/auth/signup", h.signupPost)

	return mux
}
