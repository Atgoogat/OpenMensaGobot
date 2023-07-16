package domain

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Atgoogat/openmensarobot/db"
	"github.com/Atgoogat/openmensarobot/openmensa"
	"github.com/Atgoogat/openmensarobot/telegrambotapi"
)

type Api struct {
	openmensaApi openmensa.OpenmensaApi
	scheduler    *SubscriberScheduler
	repo         db.SubscriberRepository
	tapi         *telegrambotapi.TelegramBotApi
}

func NewApi(openmensaApi openmensa.OpenmensaApi,
	repo db.SubscriberRepository,
	telegramApi *telegrambotapi.TelegramBotApi) Api {
	scheduler := NewSubscriberScheduler(nil)

	return Api{
		openmensaApi: openmensaApi,
		scheduler:    scheduler,
		repo:         repo,
		tapi:         telegramApi,
	}
}

func (api Api) Start(ctx context.Context) <-chan error {
	messages := api.tapi.GetMessageChan()
	done := make(chan error)

	go func() {
		running := true
		for running {
			select {
			case m := <-messages:
				api.processMessage(m)
			case <-ctx.Done():
				running = false
				done <- ctx.Err()
			}
		}
	}()

	return done
}

func (api Api) processMessage(msg telegrambotapi.TelegramMessage) {
	if len(msg.Text) >= 100 {
		api.tapi.SendMessage(msg.ChatID, "Your message is too long")
		return
	}
	splits := strings.Split(msg.Text, " ")
	if len(splits) == 0 {
		// should not happen as the message is guranteed to be an comman (contains at least a slash)
		return
	}

	switch splits[0] {
	case "/today":
		api.processToday(msg)
	default:
		api.tapi.SendMessage(msg.ChatID, "I could not understand your command")
	}
}

func (api Api) processToday(msg telegrambotapi.TelegramMessage) {
	text := strings.Trim(msg.Text, " \n")
	splits := strings.Split(text, " ")
	if len(splits) != 2 {
		api.tapi.SendMessage(msg.ChatID, "Expected one argument.\nExample: /today 1719")
		return
	}

	mensaID, err := strconv.ParseInt(splits[1], 10, 64)
	if err != nil {
		api.tapi.SendMessage(msg.ChatID, "Invalid Mensa ID.\nExample: /today 1719")
		return
	}

	meals, err := api.openmensaApi.ListMealsForADay(int(mensaID), time.Now(), 1, 20)
	if err != nil {
		api.tapi.SendMessage(msg.ChatID, fmt.Sprintf("Could not get Meals from Mensa %d", mensaID))
		return
	}

	if len(meals) == 0 {
		api.tapi.SendMessage(msg.ChatID, "There are no meals for today!")
		return
	}

	returnMsg := MealsToMsg(meals, openmensa.PRICE_STUDENT)
	api.tapi.SendMessage(msg.ChatID, returnMsg)
}

func (api Api) registerSubscriber(chatID, mensaID int, pushHours, pushMinutes uint) error {
	sub, err := api.repo.CreateSubscriber(chatID, mensaID, pushHours, pushMinutes)
	if err != nil {
		return err
	}

	err = api.scheduler.InsertJob(sub)
	return err
}

func (api Api) removeSubscriber(chatID int) error {
	sub, err := api.repo.FindSubscriberByChatID(chatID)
	if err != nil {
		return err
	}

	err = api.repo.DeleteSubscriberByID(sub.ID)
	if err != nil {
		return err
	}
	api.scheduler.RemoveJob(sub.ID)
	return nil
}
