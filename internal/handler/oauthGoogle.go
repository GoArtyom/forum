package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"forum/internal/models"
	"forum/pkg"
)

const (
	googleAuthURL     = "https://accounts.google.com/o/oauth2/auth"
	googleTokenURL    = "https://oauth2.googleapis.com/token"
	googleUserInfoURL = "https://www.googleapis.com/oauth2/v3/userinfo"
)

func (h *Handler) signinGoogle(w http.ResponseWriter, r *http.Request) {
	state := pkg.RandString(111)
	pkg.SetStateCookie(w, state)
	fmt.Println(h.gooleConfig.RedirectURL)
	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&state=%s&scope=profile email",
		googleAuthURL, h.gooleConfig.ClientID, h.gooleConfig.RedirectURL, state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect) // 307
}

func (h *Handler) callbackGoogle(w http.ResponseWriter, r *http.Request) {
	fmt.Println(h.gooleConfig.RedirectURL)
	state, err := r.Cookie("state")
	if err != nil {
		log.Printf("callbackGoogle: state not found: %s\n", err.Error())
		http.Error(w, "state not found", http.StatusBadRequest) // 400
		return
	}
	if r.URL.Query().Get("state") != state.Value {
		log.Println("callbackGoogle: state did not match")
		http.Error(w, "state did not match", http.StatusBadRequest) // 400
		return
	}

	
	code := r.URL.Query().Get("code")
	form := strings.NewReader(fmt.Sprintf("code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code",
		code, h.gooleConfig.ClientID, h.gooleConfig.ClientSecret, h.gooleConfig.RedirectURL))

	resp, err := http.Post(googleTokenURL, "application/x-www-form-urlencoded", form)
	if err != nil {
		log.Printf("callbackGoogle: failed to POST request:%s\n", err.Error())
		http.Error(w, "failed to request in API google", http.StatusInternalServerError) // 500
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("callbackGoogle: failed to read response body: %s\n", err.Error())
		http.Error(w, "failed to read response body", http.StatusInternalServerError) // 500
		return
	}
	access_token := getValueFromBody(body, "access_token")
	if access_token == "" {
		log.Println("callbackGoogle: access_token not found")
		http.Error(w, "failed to extract access token", http.StatusInternalServerError) // 500
		return
	}
	fmt.Println("access_token:", access_token)

	url := fmt.Sprintf("%s?access_token=%s", googleUserInfoURL, access_token)
	resp, err = http.Get(url)

	if err != nil {
		log.Printf("callbackGoogle: failed to GET request:%s\n", err.Error())
		http.Error(w, "failed to request in API google", http.StatusInternalServerError) // 500
		return
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("callbackGoogle: failed to read response body: %s\n", err.Error())
		http.Error(w, "failed to read response body", http.StatusInternalServerError) // 500
		return
	}
	fmt.Println(string(body))
	userInfo := &models.UserInfoOAuth{}
	err = getUserInfo(body, userInfo)
	if err != nil {
		log.Printf("callbackGoogle: failed to create user: %s\n", err.Error())
		http.Error(w, "failed to create user", http.StatusInternalServerError) // 500
		return
	}
	_, err = h.service.GetUserByEmail(userInfo.Email)
	if err != nil {
		// create
		createUser := &models.CreateUser{
			Name:     userInfo.Name,
			Email:    userInfo.Email,
			Password: userInfo.Sub,
		}
		err := h.service.CreateUser(createUser)
		if err != nil {
			log.Printf("callbackGoogle:CreateUser:%s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
			return
		}
	}
	signInUser := &models.SignInUser{
		Email:    userInfo.Email,
		Password: userInfo.Sub,
	}

	userId, err := h.service.SignInUser(signInUser)
	if err != nil {
		log.Printf("callbackGoogle:SignInUser:%s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}
	session, err := h.service.CreateSession(userId)
	if err != nil {
		log.Printf("callbackGoogle:CreateSession:%s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}
	pkg.SetCookie(w, session.UUID, session.ExpireAt)

	http.Redirect(w, r, "/", http.StatusSeeOther) // 303
}
