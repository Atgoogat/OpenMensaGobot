package telegrambotapi

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBotApi struct {
	api            *tgbotapi.BotAPI
	messageChannel chan TelegramMessage
}

type TelegramMessage struct {
	ChatID int64
	Text   string
}

func NewTelegramBotApi(apiToken string) (*TelegramBotApi, error) {
	api, err := tgbotapi.NewBotAPI(apiToken)
	return &TelegramBotApi{
		api: api,
	}, err
}

func (api TelegramBotApi) GetMessageChan() <-chan TelegramMessage {
	if api.messageChannel == nil {
		api.messageChannel = make(chan TelegramMessage, 10)
		go func() {
			u := tgbotapi.NewUpdate(0)
			u.Timeout = 60

			updates := api.api.GetUpdatesChan(u)

			for update := range updates {
				msg := update.Message
				if msg != nil && msg.IsCommand() {
					tmsg := TelegramMessage{
						ChatID: msg.Chat.ID,
						Text:   msg.Text,
					}
					api.messageChannel <- tmsg
				}
			}
		}()
	}
	return api.messageChannel
}

func (api TelegramBotApi) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := api.api.Send(msg)
	if err != nil {
		log.Printf("errors sending message to %d: %v", chatID, err)
	}
	return err
}