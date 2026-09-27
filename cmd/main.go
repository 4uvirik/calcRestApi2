package main

import (
	"github.com/4uvirik/calcRestApi2/config"
	"github.com/4uvirik/calcRestApi2/pkg/logger"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"os"
)

func main() {

	loadEnv()

	cfg := initConfig()

	logger := initLogger(cfg)

	logger.Info("ЗАГЛУШКА")
}

// loadEnv - загрузка переменных окружения.
func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println(".env file not found")
	}
}

// initConfig - загрузка конфига и проверка на валидацию.
func initConfig() *config.Config {
	yamlPath := os.Getenv("YAML_PATH")

	cfg, err := config.LoadConfig(yamlPath)
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
	logger, err := logger.ConfigLogger(cfg.Logger.Level)
	if err != nil {
		log.Fatalf("error creating logger: %v", err)
	}

	logger.Debug("configuration reading success", slog.Any("cfg", cfg))

	return logger
}
