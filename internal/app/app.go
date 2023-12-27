package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"forum/config"
	"forum/internal/handler"
	"forum/internal/render"
	repo "forum/internal/repository"
	"forum/internal/server"
	"forum/internal/service"
)

func RunServer(cfg *config.Config) {
	db, err := repo.NewSqliteDB(cfg)
	if err != nil {
		log.Fatalf("[ERROR]:failed to initialize db: %s\n", err.Error())
	}
	err = repo.CreateTable(db, cfg.Migrate)
	if err != nil {
		log.Fatalf("[ERROR]:failed creation table: %s\n", err.Error())
	}
	repo := repo.NewRepository(db)
	service := service.NewService(repo)
	tpl, err := render.NewTemplate()
	if err != nil {
		log.Fatalf("[ERROR]:failed to parse templates: %s\n", err.Error())
	}
	handlers := handler.NewHandler(service, tpl, cfg.GoogleConfig, cfg.GithubConfig)
	srv := new(server.Server)

	go func() {
		if err := srv.Run(cfg, handlers.InitRouters()); err != nil {
			log.Printf("[ERROR]:occured while running http server: %s\n", err.Error())
		}
	}()

	go handler.CleanupVisitors()

	log.Println("[OK]:listening on: http://localhost" + cfg.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("[ERROR]:occured on server shutting down : %s", err.Error())
	} else {
		log.Println("[OK]:server shutdown was successful")
	}

	if err := db.Close(); err != nil {
		log.Printf("[ERROR]:occured on db connection close : %s", err.Error())
	} else {
		log.Println("[OK]:db shutdown was successful")
	}
}
