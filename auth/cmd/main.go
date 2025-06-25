package main

import (
	"auth"
	"auth/pkg/config"
	"auth/pkg/handler"
	"auth/pkg/logger"
	"auth/pkg/repository"
	"auth/pkg/service"
	"context"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	cfg, err := config.LoadConfig("./configs/config.yml")
	if err != nil {
		log.Fatal(err)
	}

	logger.Init(cfg.Logger)

	logrus.Debug("Config loaded")

	db, err := repository.NewPostgresDB(cfg.DataSource)
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	logrus.Debug("DB connected")

	repos := repository.NewRepository(db)
	services := service.NewService(repos, cfg.Auth)
	handlers := handler.NewHandler(services)

	srv := new(auth.Server)
	go func() {
		if err := srv.Run(cfg.Main.Port, handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()
	logrus.Info("Auth server started on port: ", cfg.Main.Port)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logrus.Println("Server shutting down")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}
