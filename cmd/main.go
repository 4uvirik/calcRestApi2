package main

import (
	"github.com/4uvirik/calcRestApi2/config"
	"github.com/4uvirik/calcRestApi2/pkg/logger"
	"log"
	"log/slog"
)

func main() {

	cfg := initConfig()

	logger := initLogger(cfg)

	logger.Info("ЗАГЛУШКА")
}

// initConfig - загрузка конфига и проверка на валидацию.
func initConfig() *config.Config {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("config validate failed: %v", err)
	}

	return cfg
}

// initLogger - загрузка и настройка логера.
func initLogger(cfg *config.Config) *slog.Logger {
	logger, err := logger.ConfigLogger(cfg.Logger.Environment)
	if err != nil {
		log.Fatalf("error creating logger: %v", err)
	}

	logger.Debug("configuration reading success", slog.Any("cfg", cfg))

	return logger
}
