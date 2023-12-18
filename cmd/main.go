package main

import (
	"flag"

	"forum/config"
	"forum/internal/app"
)

func main() {
	configPath := flag.String("cfg", "./config/config.json", "USAGE --cfg=path_to_config_file")
	flag.Parse()
	cfg := config.InitConfig(*configPath)
	app.RunServer(cfg)
}
