package handler

import "net/http"

func (h Handler) InitRouters() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", h.index)

	mux.HandleFunc("/signin", h.signinGET)
	mux.HandleFunc("/signup", h.signupGET)

	mux.HandleFunc("/auth/signin", h.signinPOST)
	mux.HandleFunc("/auth/signup", h.signupPOST)

	mux.Handle("/post/create", h.authorization(http.HandlerFunc(h.createPostPOST)))
	mux.HandleFunc("/post/", h.onePostGET)

	return h.sessionMiddleware(mux)
}
