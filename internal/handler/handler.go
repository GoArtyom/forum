package handler

import (
	"html/template"
	"log"
	"net/http"

	"forum/config"
	"forum/internal/render"
	"forum/internal/service"
)

type Handler struct {
	service      *service.Service
	template     *template.Template
	gooleConfig  config.GoogleConfig
	githubConfig config.GithubConfig
}

func NewHandler(service *service.Service, tpl *template.Template, googleCfg config.GoogleConfig, githubCfg config.GithubConfig) *Handler {
	return &Handler{
		service:      service,
		template:     tpl,
		gooleConfig:  googleCfg,
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
