package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Port string `json:"port"`
	Migrate string `json:"migrate"`
	DB   struct {
		Driver string `json:"driver"`
		DSN    string `json:"dsn"`
	}
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
