package handler

import (
	"database/sql"
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
	googleAuthURL     = "https://accounts.google.com/o/oauth2/auth"
	googleTokenURL    = "https://oauth2.googleapis.com/token"
	googleUserInfoURL = "https://www.googleapis.com/oauth2/v3/userinfo"
)

func (h *Handler) signinGoogle(w http.ResponseWriter, r *http.Request) {
	state := pkg.RandString(111)
	pkg.SetCookie(w, "state", state, time.Now().Add(time.Minute))

	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&state=%s&scope=profile email",
		googleAuthURL, h.googleConfig.ClientID, h.googleConfig.RedirectURL, state)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect) // 307
}

func (h *Handler) callbackGoogle(w http.ResponseWriter, r *http.Request) {
	state, err := pkg.GetCookie(r, "state")
	if err != nil {
		log.Printf("callbackGoogle: state not found: %s\n", err.Error())
		h.renderError(w, http.StatusBadRequest) // 400
		return
	}

	pkg.DeleteCookie(w, "state")

	if r.URL.Query().Get("state") != state.Value {
		log.Println("callbackGoogle: state did not match")
		h.renderError(w, http.StatusBadRequest) // 400
		return
	}

	code := r.URL.Query().Get("code")

	form := strings.NewReader(fmt.Sprintf("code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code",
		code, h.googleConfig.ClientID, h.googleConfig.ClientSecret, h.googleConfig.RedirectURL))

	resp, err := http.Post(googleTokenURL, "application/x-www-form-urlencoded", form)
	if err != nil {
		log.Printf("callbackGoogle: failed to POST request:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("callbackGoogle: failed to read response body: %s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	access_token := getValueFromBody(body, "access_token")
	if access_token == "" {
		log.Println("callbackGoogle: access_token not found")
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	userInfo, err := h.getUserInfoFromApi(access_token, googleUserInfoURL)
	if err != nil {
		log.Printf("callbackgithub: failed to GET request:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	userInfoGoogle := &models.UserInfoGoogle{}

	err = json.Unmarshal(userInfo, &userInfoGoogle)
	if err != nil {
		log.Printf("callbackgithub: failed unmarshal userInfo: %s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}
	_, err = h.service.GetUserByEmail(userInfoGoogle.Email)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("callbackgithub: failed GetUserByEmail: %s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	if err == sql.ErrNoRows {
		createUser := &models.CreateUser{
			Name:     userInfoGoogle.Name,
			Email:    userInfoGoogle.Email,
			Password: userInfoGoogle.Sub,
			Mode:     models.GoogleMode,
		}

		err := h.service.CreateUser(createUser)
		if err != nil {
			log.Printf("callbackGoogle:CreateUser:%s\n", err.Error())
			h.renderError(w, http.StatusInternalServerError) // 500
			return
		}

	}

	signInUser := &models.SignInUser{
		Name:     userInfoGoogle.Name,
		Email:    userInfoGoogle.Email,
		Password: userInfoGoogle.Sub,
		Mode:     models.GoogleMode,
	}

	userId, err := h.service.SignInUser(signInUser)
	if err != nil {
		log.Printf("callbackGoogle:SignInUser:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	session, err := h.service.CreateSession(userId)
	if err != nil {
		log.Printf("callbackGoogle:CreateSession:%s\n", err.Error())
		h.renderError(w, http.StatusInternalServerError) // 500
		return
	}

	pkg.SetCookie(w, "UUID", session.UUID, session.ExpireAt)

	http.Redirect(w, r, "/", http.StatusSeeOther) // 303
}
