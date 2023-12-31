package telegrambotapi

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	DESCRIPTION = "This bot can send you menus from openmensa.org. Look at their site and copy the mensaID from the url."
)

type TelegramBotApi struct {
	api            *tgbotapi.BotAPI
	messageChannel chan TelegramMessage
}

type TelegramMessage struct {
	ChatID int64
	// Text is trimmed
	Text string
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
					text := strings.Trim(msg.Text, " ")
					if text == "/help" {
						err := api.sendHelpMsg(msg.Chat.ID)
						if err != nil {
							log.Printf("error sending help message: %v", err)
						}
						return
					}
					tmsg := TelegramMessage{
						ChatID: msg.Chat.ID,
						Text:   text,
					}
					api.messageChannel <- tmsg
				}
			}
		}()
	}
	return api.messageChannel
}

func (api TelegramBotApi) SendHtmlMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "html"
	_, err := api.api.Send(msg)
	if err != nil {
		log.Printf("errors sending message to %d: %v", chatID, err)
	}
	return err
}

func (api TelegramBotApi) sendHelpMsg(chatID int64) error {
	cmds, err := api.api.GetMyCommands()
	if err != nil {
		return err
	}

	msg := []string{DESCRIPTION, ""}
	for _, cmd := range cmds {
		msg = append(msg, cmd.Command+" - "+cmd.Description)
	}
	return api.SendHtmlMessage(chatID, strings.Join(msg, "\n"))
}
