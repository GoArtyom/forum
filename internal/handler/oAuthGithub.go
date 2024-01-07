package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"forum/internal/models"
	"forum/pkg"
)

const (
	githubAuthURL     = "https://github.com/login/oauth/authorize"
	githubTokenURL    = "https://github.com/login/oauth/access_token"
	githubUserInfoURL = "https://api.github.com/user"
)

func (h *Handler) signinGithub(w http.ResponseWriter, r *http.Request) {
	state := pkg.RandString(111)

	pkg.SetCookie(w, "state", state, time.Now().Add(time.Minute*10))

	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&&response_type=code&allow_signup=true&state=%s&scope=user:email",
		githubAuthURL, h.githubConfig.ClientID, h.githubConfig.RedirectURL, state)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect) // 307
}

func (h *Handler) callbackGithub(w http.ResponseWriter, r *http.Request) {
	state, err := pkg.GetCookie(r, "state")
	if err != nil {
		log.Printf("callbackgithub: state not found: %s\n", err.Error())
		h.renderError(w, http.StatusBadRequest) // 400
		return
	}
	pkg.DeleteCookie(w, "state")

	if r.URL.Query().Get("state") != state.Value {
		log.Println("callbackgithub: state did not match")
		h.renderError(w, http.StatusBadRequest) // 400
		return
	}

	code := r.URL.Query().Get("code")

	form := strings.NewReader(fmt.Sprintf("code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code",
		code, h.githubConfig.ClientID, h.githubConfig.ClientSecret, h.githubConfig.RedirectURL))

	resp, err := http.Post(githubTokenURL, "application/x-www-form-urlencoded", form)
	if err != nil {
		log.Printf("callbackgithub: failed to POST request:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("callbackgithub: failed to read response body: %s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}
	access_token := getValueFromURL(string(body), "access_token")
	if access_token == "" {
		log.Println("callbackgithub: access_token not found")
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	userInfo, err := h.getUserInfoFromApi(access_token, githubUserInfoURL)
	if err != nil {
		log.Printf("callbackgithub: failed to GET request:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	userInfoGitHub := &models.UserInfoGitHub{}

	err = json.Unmarshal(userInfo, &userInfoGitHub)
	if err != nil {
		log.Printf("callbackgithub: failed unmarshal userInfo: %s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	_, err = h.service.GetByEmail(userInfoGitHub.NodeId)
	if err != nil {
		createUser := &models.CreateUser{
			Name:     userInfoGitHub.Login,
			Email:    userInfoGitHub.NodeId,
			Password: userInfoGitHub.NodeId,
			Mode:     models.GitHubMode,
		}

		err := h.service.User.Create(createUser)
		if err != nil {
			log.Printf("callbackgithub:CreateUser:%s\n", err.Error())
			h.renderError(w, http.StatusInternalServerError) // 500
			return
		}

	}

	signInUser := &models.SignInUser{
		Name:     userInfoGitHub.Login,
		Email:    userInfoGitHub.NodeId,
		Password: userInfoGitHub.NodeId,
		Mode:     models.GitHubMode,
	}

	userId, err := h.service.User.SignIn(signInUser)
	if err != nil {
		log.Printf("callbackgithub:SignInUser:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	session, err := h.service.Session.Create(userId)
	if err != nil {
		log.Printf("callbackgithub:Create:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	pkg.SetCookie(w, "UUID", session.UUID, session.ExpireAt)

	http.Redirect(w, r, "/", http.StatusSeeOther) // 303
}
