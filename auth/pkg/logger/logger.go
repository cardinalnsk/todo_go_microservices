package logger

import (
	"auth/pkg/config"
	"os"

	"github.com/sirupsen/logrus"
)

func Init(cfg config.LoggerConfig) {
	if cfg.Level == "" {
		cfg.Level = "info"
	}

	logrus.SetOutput(os.Stdout)

	parsedLevel, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		logrus.Fatalf("invalid log cfg: %s", cfg.Level)
	}
	logrus.SetLevel(parsedLevel)

	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		PrettyPrint:     cfg.Pretty,
	})
}
