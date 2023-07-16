package config

import (
	"log"
	"os"

	"github.com/Atgoogat/openmensarobot/telegrambotapi"
)

func NewTelegramBotApi() *telegrambotapi.TelegramBotApi {
	token, ok := os.LookupEnv(TELEGRAM_TOKEN)
	if !ok {
		log.Fatalf("Env var %s not found", TELEGRAM_TOKEN)
	}

	api, err := telegrambotapi.NewTelegramBotApi(token)
	if err != nil {
		log.Fatalf("error while setting up telegram api: %v\n", err)
	}
	log.Println("connected to telegram api")
	return api
}
