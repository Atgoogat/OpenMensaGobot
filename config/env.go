package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	DB_URL         = "OPENMENSAROBOT_DB_URL"
	TELEGRAM_TOKEN = "OPENMENSAROBOT_TELEGRAM_TOKEN"
)

func LoadEnv() error {
	return godotenv.Load(".env")
}

func Getenv(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return value, fmt.Errorf("environement key %s not found", key)
	}
	return value, nil
}
