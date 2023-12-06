package main

import (
	"flag"

	"forum/config"
	"forum/internal/app"
)

func main() {
	configPath := flag.String("cfg", "./config/config.json", "----")
	flag.Parse()
	cfg := config.InitConfig(*configPath)
	app.RunServer(cfg)
}
