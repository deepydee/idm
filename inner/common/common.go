package common

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

// Config общая конфигурация всего приложения
type Config struct {
	DbDriverName string `validate:"required"`
	Dsn          string `validate:"required"`
}

// GetConfig получение конфигурации из .env файла или переменных окружения
func GetConfig(envFile string) Config {
	_ = godotenv.Load(envFile)

	dsn := fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_CONNECTION"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	return Config{
		DbDriverName: os.Getenv("DB_CONNECTION"),
		Dsn:          dsn,
	}
}
