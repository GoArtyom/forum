package handler

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"forum/config"
	"forum/internal/render"
	"forum/internal/service"
)

type Handler struct {
	service      *service.Service
	template     *template.Template
	googleConfig config.GoogleConfig
	githubConfig config.GithubConfig
}

func NewHandler(service *service.Service, tpl *template.Template, googleCfg config.GoogleConfig, githubCfg config.GithubConfig) *Handler {
	return &Handler{
		service:      service,
		template:     tpl,
		googleConfig: googleCfg,
		githubConfig: githubCfg,
	}
}

func (h *Handler) renderPage(w http.ResponseWriter, file string, data *render.Data) {
	err := h.template.ExecuteTemplate(w, file, data)
	if err != nil {
		log.Printf("ExecuteTemplate:%s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
	}
}

func (h *Handler) renderError(w http.ResponseWriter, code int) {
	w.WriteHeader(code)

	err := h.template.ExecuteTemplate(w, "error.html", struct {
		Code int
		Text string
	}{
		Code: code,
		Text: http.StatusText(code),
	})
	if err != nil {
		log.Printf("ExecuteTemplate:%s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
	}
}

func (h *Handler) getUserInfoFromApi(accessToken string, userInfoURL string) ([]byte, error) {
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
