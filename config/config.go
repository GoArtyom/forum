package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Port    string `json:"port"  env-default:":8080"`
	Migrate string `json:"migrate"`
	DB      struct {
		Driver string `json:"driver"`
		DSN    string `json:"dsn"`
	}
	GoogleConfig GoogleConfig `json:"google_config"`
	GithubConfig GithubConfig `json:"github_config"`
}

type GoogleConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_url"`
}

type GithubConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_url"`
}

func InitConfig(path string) *Config {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("error open config file: %s", err.Error())
	}
	defer file.Close()
	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		log.Fatalf("error decoding config file: %s", err.Error())
	}
	return &config
}
