package main

import (
	apigateway "api-gateway"
	"api-gateway/pkg/config"
	"api-gateway/pkg/handler"
	"api-gateway/pkg/logger"
	"context"
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

	handlers := handler.NewHandler(cfg.Gateway, cfg.Auth, cfg.Todo)

	srv := new(apigateway.Server)
	go func() {
		if err := srv.Run(cfg.Main.Port, handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Info("Gateway service started on port: ", cfg.Main.Port)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logrus.Println("Gateway shutting down")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}
}
