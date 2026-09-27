package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App    AppConfig    `yaml:"app"`
	Logger LoggerConfig `yaml:"logger"`

	// Добавлять новые секции только по мере необходимости:
	//
	// Database DatabaseConfig `yaml:"database"`
	// Redis    RedisConfig    `yaml:"redis"`
	// Kafka    KafkaConfig    `yaml:"kafka"`
	// Storage  StorageConfig  `yaml:"storage"`
}

type AppConfig struct {
	Name string `env:"APP_NAME" yaml:"name"`

	// Host/Port нужны, если приложение поднимает HTTP/gRPC сервер.
	Host string `env:"APP_HOST" yaml:"host"`
	Port string `env:"APP_PORT" yaml:"port"`

	// Таймаут корректного завершения приложения.
	// Обычно достаточно значения по умолчанию.
	ShutdownTimeout time.Duration `env:"APP_SHUTDOWN_TIMEOUT" envDefault:"5s" yaml:"shutdown_timeout"`
}

type LoggerConfig struct {
	// Режим работы логгера:
	// local — красивый вывод для локальной разработки
	// dev — JSON + Debug
	// prod — JSON + Info
	// любое неизвестное значение — default.
	Environment string `env:"LOG_ENV" envDefault:"local" yaml:"environment"`
}

func loadFromEnv(config *Config) error {
	if _, err := os.Stat(".env"); err == nil {
		_ = godotenv.Load(".env")
	}

	if err := env.Parse(config); err != nil {
		return fmt.Errorf("failed to parse environment variables: %w", err)
	}

	return nil
}

func loadFromYaml(path string, config *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read yaml file: %w", err)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	return nil
}

// LoadConfig загружает конфигурацию.
// Приоритет:
// 1. Переменные окружения / .env
// 2. YAML-файл
//
// Важно: если появятся обязательные поля,
// после загрузки необходимо вызвать Validate().
func LoadConfig(yamlPath string) (*Config, error) {
	config := &Config{}

	if err := loadFromEnv(config); err == nil {
		return config, nil
	}

	if err := loadFromYaml(yamlPath, config); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate проверяет обязательные параметры конфигурации.
// При добавлении новых обязательных полей их необходимо
// добавить сюда.
// Например:
//
//	if cfg.Database.Host == "" {
//	    errorsMsg = append(errorsMsg, "invalid database host")
//	}
//
// Для необязательных параметров проверка обычно не требуется,
// если для них задан envDefault или есть корректное значение по умолчанию.
func (cfg *Config) Validate() error {
	var errorsMsg []string

	if cfg.App.Name == "" {
		errorsMsg = append(errorsMsg, "invalid app name")
	}

	if cfg.App.Host == "" {
		errorsMsg = append(errorsMsg, "invalid app host")
	}

	if cfg.App.Port == "" {
		errorsMsg = append(errorsMsg, "invalid app port")
	}

	if len(errorsMsg) > 0 {
		return fmt.Errorf("%s", strings.Join(errorsMsg, ": "))
	}

	return nil
}
