package domain

import (
	"context"
	"log"

	"github.com/Atgoogat/openmensarobot/telegrambotapi"
)

type MsgService struct {
	telegrambotapi *telegrambotapi.TelegramBotApi
	cliService     CliService
}

func NewMsgService(telegrambotapi *telegrambotapi.TelegramBotApi, cs CliService) MsgService {
	return MsgService{
		telegrambotapi: telegrambotapi,
		cliService:     cs,
	}
}

func (ms MsgService) StartReceivingMessages(ctx context.Context) <-chan error {
	messages := ms.telegrambotapi.GetMessageChan()

	errChannel := make(chan error, 1)

	go func() {
		running := true
		for running {
			select {
			case m := <-messages:
				err := ms.handleMessage(m)
				if err != nil {
					log.Printf("error while handling message (%#v) err (%v)", m, err)
					ms.telegrambotapi.SendHtmlMessage(m.ChatID, err.Error())
				}
			case <-ctx.Done():
				running = false
				errChannel <- ctx.Err()
			}
		}
	}()
	return errChannel
}

func (ms MsgService) handleMessage(msg telegrambotapi.TelegramMessage) error {
	text, err := ms.cliService.ParseAndExecuteCommand(msg)
	if err != nil {
		return err
	}
	return ms.telegrambotapi.SendHtmlMessage(msg.ChatID, text)
}
