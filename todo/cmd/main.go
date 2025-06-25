package main

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	"syscall"
	"todo"
	"todo/pkg/config"
	"todo/pkg/handler"
	"todo/pkg/logger"
	"todo/pkg/repository"
	"todo/pkg/service"
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
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	go func() {
		if err := srv.Run(cfg.Main.Port, handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Info("Todo server started on port: ", cfg.Main.Port)

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
