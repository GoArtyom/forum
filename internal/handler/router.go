package handler

import "net/http"

func (h *Handler) InitRouters() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", h.index)

	mux.HandleFunc("/signin", h.signinGET)
	mux.HandleFunc("/auth/signin", h.signinPOST)

	mux.HandleFunc("/signup", h.signupGET)
	mux.HandleFunc("/auth/signup", h.signupPOST)

	mux.HandleFunc("/post/", h.onePostGET)
	mux.Handle("/post/create", h.authorization(http.HandlerFunc(h.createPostPOST)))

	mux.Handle("/comment/create", h.authorization(http.HandlerFunc(h.createCommentPOST)))

	mux.Handle("/myposts", h.authorization(http.HandlerFunc(h.myPostsGET)))
	mux.Handle("/auth/signout", h.authorization(http.HandlerFunc(h.signoutPOST)))

	return h.sessionMiddleware(mux)
}
