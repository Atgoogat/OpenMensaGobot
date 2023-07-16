package config

import (
	"github.com/joho/godotenv"
)

const (
	DB_URL         = "OPENMENSAROBOT_DB_URL"
	TELEGRAM_TOKEN = "OPENMENSAROBOT_TELEGRAM_TOKEN"
)

func LoadEnv() error {
	return godotenv.Load(".env")
}
