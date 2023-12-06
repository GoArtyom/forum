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
	db, err := repo.NewSqliteDB(&repo.Config{
		Driver: cfg.DB.Driver,
		Dsn:    cfg.DB.DSN,
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	repo := repo.NewRepository(db)
	service := service.NewService(repo)
	tpl, err := render.NewTemplate()
	if err != nil {
		log.Fatalf("failed to parse templates: %s", err.Error())
	}
	handler := handler.NewHandler(service, tpl)
	svr := new(server.Server)
	go func() {
		if err := svr.Run(cfg, handler.InitRouters()); err != nil {
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := svr.Shutdown(context.Background()); err != nil {
		log.Printf("error occured on server shutting down : %s", err.Error())
	}
	if err := db.Close(); err != nil {
		log.Printf("error occured on db connection close : %s", err.Error())
	}
}
