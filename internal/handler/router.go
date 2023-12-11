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
	
	mux.Handle("/auth/signout", h.authorization(http.HandlerFunc(h.signoutPOST)))
	
	mux.HandleFunc("/post/", h.onePostGET)
	mux.Handle("/post/create", h.authorization(http.HandlerFunc(h.createPostPOST)))
	
	mux.Handle("/comment/create", h.authorization(http.HandlerFunc(h.createCommentPOST)))
	
	mux.HandleFunc("/filterposts", h.filterPostsGET)
	mux.Handle("/myposts", h.authorization(http.HandlerFunc(h.myPostsGET)))
	mux.Handle("/likeposts", h.authorization(http.HandlerFunc(h.likePostsGET)))
	
	return h.sessionMiddleware(mux)
}
