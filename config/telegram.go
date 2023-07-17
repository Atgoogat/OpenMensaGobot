package config

import (
	"log"
	"sync"

	"github.com/Atgoogat/openmensarobot/telegrambotapi"
)

var (
	telegramApi     *telegrambotapi.TelegramBotApi
	telegramApiOnce sync.Once
)

func NewTelegramBotApi() (*telegrambotapi.TelegramBotApi, error) {
	var e error
	telegramApiOnce.Do(func() {
		log.Println("connecting to telegram api")
		token, err := Getenv(TELEGRAM_TOKEN)
		if err != nil {
			e = err
			return
		}

		api, err := telegrambotapi.NewTelegramBotApi(token)
		if err != nil {
			e = err
			return
		}
		log.Println("connected to telegram api")
		telegramApi = api
	})
	return telegramApi, e
}
