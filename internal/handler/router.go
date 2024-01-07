package handler

import "net/http"

func (h *Handler) InitRouters() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", h.indexGET)

	mux.HandleFunc("/signin", h.signinGET)
	mux.HandleFunc("/auth/signin", h.signinPOST)

	mux.HandleFunc("/signup", h.signupGET)
	mux.HandleFunc("/auth/signup", h.signupPOST)

	mux.HandleFunc("/auth/google/signin", h.signinGoogle)
	mux.HandleFunc("/google/callback", h.callbackGoogle)

	mux.HandleFunc("/auth/github/signin", h.signinGithub)
	mux.HandleFunc("/github/callback", h.callbackGithub)

	mux.Handle("/auth/signout", h.authUser(http.HandlerFunc(h.signoutPOST)))

	mux.HandleFunc("/post", h.onePostGET)
	mux.Handle("/post/create", h.authUser(http.HandlerFunc(h.createPostGET_POST)))
	mux.Handle("/post/delete", h.authUser(http.HandlerFunc(h.deletePostDELETE))) //++
	mux.Handle("/post/update", h.authUser(http.HandlerFunc(h.updatePostGET_POST)))

	mux.Handle("/comment/create", h.authUser(http.HandlerFunc(h.createCommentPOST)))
	mux.Handle("/comment/delete", h.authUser(http.HandlerFunc(h.deleteCommentDELETE))) //++
	mux.Handle("/comment/update", h.authUser(http.HandlerFunc(h.updateCommentGET_POST)))

	mux.Handle("/post/vote/create", h.authUser(http.HandlerFunc(h.createPostVotePOST)))
	mux.Handle("/comment/vote/create", h.authUser(http.HandlerFunc(h.createCommentVotePOST)))

	mux.HandleFunc("/filterposts", h.filterPostsGET)
	mux.Handle("/myactivity", h.authUser(http.HandlerFunc(h.myActivityGET)))
	mux.Handle("/mynotifications", h.authUser(http.HandlerFunc(h.myNotificationsGET)))

	// moderation
	mux.Handle("/moderator/request", h.authUser(http.HandlerFunc(h.moderatorRequestPATCH)))           //++
	mux.Handle("/post/reporting", h.authUser(h.authModerator(http.HandlerFunc(h.reportingPostPOST)))) //++

	// admin
	mux.Handle("/admin", h.authUser(h.authAdmin(http.HandlerFunc(h.adminGET))))                                     //++
	mux.Handle("/admin/report", h.authUser(h.authAdmin(http.HandlerFunc(h.adminReportDELETE))))                     //++
	mux.Handle("/admin/categories/delete", h.authUser(h.authAdmin(http.HandlerFunc(h.adminCategoriesDELETE))))      //++
	mux.Handle("/admin/categories/create", h.authUser(h.authAdmin(http.HandlerFunc(h.adminCategoriesCREATE))))      //++
	mux.Handle("/admin/moderator-request", h.authUser(h.authAdmin(http.HandlerFunc(h.adminModeratorRequestPATCH)))) //++

	return h.recoverPanic(h.secureHeaders(h.sessionMiddleware(h.limit(5, 5, mux))))
}
