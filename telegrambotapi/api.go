package telegrambotapi

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	DESCRIPTION = "This bot can send you menus from openmensa.org. Look at their site and copy the mensaID from the url."

	// Telegram rejects messages longer than 4096 characters
	MaxMessageLength = 4096
)

type TelegramBotApi struct {
	api            *tgbotapi.BotAPI
	messageChannel chan TelegramMessage
	once           sync.Once
}

type TelegramMessage struct {
	ChatID int64
	// Text is trimmed and normalized to "/command args" (a "@botname" suffix is removed)
	Text string
}

func NewTelegramBotApi(apiToken string) (*TelegramBotApi, error) {
	api, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		return nil, err
	}
	return &TelegramBotApi{
		api: api,
	}, nil
}

func (api *TelegramBotApi) GetMessageChan() <-chan TelegramMessage {
	api.once.Do(func() {
		api.messageChannel = make(chan TelegramMessage, 10)
		go func() {
			u := tgbotapi.NewUpdate(0)
			u.Timeout = 60

			updates := api.api.GetUpdatesChan(u)

			for update := range updates {
				msg := update.Message
				if msg == nil || !msg.IsCommand() {
					continue
				}

				// normalize "/cmd@BotName args" to "/cmd args"
				text := "/" + msg.Command()
				if args := strings.TrimSpace(msg.CommandArguments()); args != "" {
					text += " " + args
				}

				if text == "/help" {
					if err := api.sendHelpMsg(msg.Chat.ID); err != nil {
						log.Printf("error sending help message: %v", err)
					}
					continue
				}
				api.messageChannel <- TelegramMessage{
					ChatID: msg.Chat.ID,
					Text:   text,
				}
			}
		}()
	})
	return api.messageChannel
}

// IsChatUnreachable reports whether the error means the bot can no longer
// write to the chat (bot blocked, kicked, chat deleted, user deactivated).
func IsChatUnreachable(err error) bool {
	var tgErr *tgbotapi.Error
	if errors.As(err, &tgErr) {
		return tgErr.Code == http.StatusForbidden
	}
	return false
}

func (api *TelegramBotApi) SendHtmlMessage(chatID int64, text string) error {
	for _, part := range splitMessage(text, MaxMessageLength) {
		msg := tgbotapi.NewMessage(chatID, part)
		msg.ParseMode = tgbotapi.ModeHTML
		if _, err := api.api.Send(msg); err != nil {
			log.Printf("errors sending message to %d: %v", chatID, err)
			return err
		}
	}
	return nil
}

// splitMessage splits text into chunks of at most limit runes,
// preferring to cut at blank lines, then at line breaks.
func splitMessage(text string, limit int) []string {
	var parts []string
	for utf8.RuneCountInString(text) > limit {
		cut := runeOffset(text, limit)
		head := text[:cut]
		if i := strings.LastIndex(head, "\n\n"); i > 0 {
			cut = i
		} else if i := strings.LastIndex(head, "\n"); i > 0 {
			cut = i
		}
		parts = append(parts, text[:cut])
		text = strings.TrimLeft(text[cut:], "\n")
	}
	return append(parts, text)
}

func runeOffset(s string, n int) int {
	i := 0
	for pos := range s {
		if i == n {
			return pos
		}
		i++
	}
	return len(s)
}

func (api *TelegramBotApi) sendHelpMsg(chatID int64) error {
	cmds, err := api.api.GetMyCommands()
	if err != nil {
		return err
	}

	msg := []string{DESCRIPTION, ""}
	for _, cmd := range cmds {
		msg = append(msg, "/"+cmd.Command+" - "+cmd.Description)
	}
	return api.SendHtmlMessage(chatID, strings.Join(msg, "\n"))
}
