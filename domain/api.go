package domain

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Atgoogat/openmensarobot/db"
	"github.com/Atgoogat/openmensarobot/openmensa"
	"github.com/Atgoogat/openmensarobot/telegrambotapi"
)

var ErrSubcriptionAlreadyExists = errors.New("this subscription already exists")

type Api struct {
	openmensaApi openmensa.OpenmensaApi
	scheduler    *SubscriberScheduler
	repo         db.SubscriberRepository
	tapi         *telegrambotapi.TelegramBotApi
	formatter    []TextFormatter
}

func NewApi(openmensaApi openmensa.OpenmensaApi,
	repo db.SubscriberRepository,
	telegramApi *telegrambotapi.TelegramBotApi) Api {

	api := Api{
		openmensaApi: openmensaApi,
		repo:         repo,
		tapi:         telegramApi,
	}

	scheduler := NewSubscriberScheduler(func(subID uint) error {
		sub, err := repo.FindSubscriberById(subID)
		if err != nil {
			return err
		}

		api.sendTodaysPlan(sub.ChatID, sub.MensaID)
		return nil
	})
	api.scheduler = scheduler
	return api
}

func (api *Api) SetFormatter(formater ...TextFormatter) {
	api.formatter = formater
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

func (api Api) sendErrorMsg(err error, chatID int64, format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	log.Printf("%s: %v\n", msg, err)
	_ = api.tapi.SendMessage(chatID, msg)
}

func (api Api) processMessage(msg telegrambotapi.TelegramMessage) {
	if len(msg.Text) >= 100 {
		api.sendErrorMsg(nil, msg.ChatID, "Your message is too long")
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
	case "/subscribe":
		api.processSubscribe(msg)
	case "/unsubscribe":
		api.processUnsubscribe(msg)
	default:
		api.sendErrorMsg(nil, msg.ChatID, "I could not understand your command")
	}
}

func (api Api) processToday(msg telegrambotapi.TelegramMessage) {
	text := strings.Trim(msg.Text, " ")
	splits := strings.Split(text, " ")
	if len(splits) != 2 {
		api.sendErrorMsg(nil, msg.ChatID, "Expected one argument.\nExample: /today 1719")
		return
	}

	mensaID, err := strconv.ParseInt(splits[1], 10, 64)
	if err != nil {
		api.sendErrorMsg(err, msg.ChatID, "Invalid Mensa ID.\nExample: /today 1719")
		return
	}

	api.sendTodaysPlan(msg.ChatID, int(mensaID))
}

// Format "/subscribe [mensaID] (HH:MM)"
func (api Api) processSubscribe(msg telegrambotapi.TelegramMessage) {
	text := strings.Trim(msg.Text, " ")
	splits := strings.Split(text, " ")

	if len(splits) != 2 && len(splits) != 3 {
		api.sendErrorMsg(nil, msg.ChatID, "Unexpected format.\nExample: /subscribe [mensaID] (HH:MM)")
		return
	}

	mensaID, err := strconv.ParseInt(splits[1], 10, 64)
	if err != nil {
		api.sendErrorMsg(err, msg.ChatID, "Invalid Mensa ID.\nExample: /subscribe 1719 09:00")
		return
	}

	var pushHours uint = 9
	var pushMinutes uint = 0
	if len(splits) == 3 {
		n, err := fmt.Sscanf(splits[2], "%02d:%02d", &pushHours, &pushMinutes)
		if err != nil || n != 2 {
			api.sendErrorMsg(err, msg.ChatID, "Could not parse push time.\nExample: /subscribe 1719 09:00")
			return
		}
	}

	if pushHours > 23 || pushMinutes > 59 {
		api.sendErrorMsg(nil, msg.ChatID, "Hours must be less than 23 and Minutes less than 60")
		return
	}

	err = api.registerSubscriber(msg.ChatID, int(mensaID), pushHours, pushMinutes)
	if err != nil {
		if err == ErrSubcriptionAlreadyExists {
			api.sendErrorMsg(err, msg.ChatID, "This subscription already exists")
			return
		}
		api.sendErrorMsg(err, msg.ChatID, "Could not register subscribtion.")
		return
	}
	api.tapi.SendMessage(msg.ChatID, fmt.Sprintf("You will receive the next update on %02d:%02d", pushHours, pushMinutes))
}

// Format /unsubscribe [mensaID]
func (api Api) processUnsubscribe(msg telegrambotapi.TelegramMessage) {
	text := strings.Trim(msg.Text, " ")
	splits := strings.Split(text, " ")

	if len(splits) != 2 {
		api.sendErrorMsg(nil, msg.ChatID, "Unexpected format.\nExample: /unsubscribe [mensaID]")
		return
	}

	mensaID, err := strconv.ParseInt(splits[1], 10, 64)
	if err != nil {
		api.sendErrorMsg(err, msg.ChatID, "Invalid Mensa ID.\nExample: /unsubscribe 1719")
		return
	}

	err = api.removeSubscriber(msg.ChatID, int(mensaID))
	if err != nil {
		api.sendErrorMsg(err, msg.ChatID, "Could not unsubscribe.")
		return
	}
	api.tapi.SendMessage(msg.ChatID, fmt.Sprintf("You wont receive anymore updates for Mensa %d", mensaID))
}

// Will create or override subscription
func (api Api) registerSubscriber(chatID int64, mensaID int, pushHours, pushMinutes uint) error {
	found, err := api.repo.ExistsSubscriberWithChatIDAndMensaID(chatID, mensaID)
	if err != nil {
		return err
	}
	if found {
		return ErrSubcriptionAlreadyExists
	}

	sub, err := api.repo.CreateSubscriber(chatID, mensaID, pushHours, pushMinutes)
	if err != nil {
		return err
	}

	err = api.scheduler.InsertJob(sub)
	if err != nil {
		return err
	}
	log.Printf("Registered new subscriber for Mensa %d\n", mensaID)
	return nil
}

func (api Api) removeSubscriber(chatID int64, mensaID int) error {
	sub, err := api.repo.FindSubscriberByChatIDAndMensaID(chatID, mensaID)
	if err != nil {
		return err
	}

	err = api.repo.DeleteSubscriberByID(sub.ID)
	if err != nil {
		return err
	}
	api.scheduler.RemoveJob(sub.ID)

	log.Printf("Unsubscibred from Mensa %d\n", mensaID)
	return nil
}

func (api Api) sendTodaysPlan(chatID int64, mensaID int) {
	meals, err := api.openmensaApi.ListMealsForADay(int(mensaID), time.Now(), 1, 20)
	if err != nil {
		api.sendErrorMsg(err, chatID, fmt.Sprintf("Could not get Meals from Mensa %d", mensaID))
		return
	}

	if len(meals) == 0 {
		api.tapi.SendMessage(chatID, "There are no meals for today!")
		return
	}

	returnMsg := MealsToMsg(meals, openmensa.PRICE_STUDENT)
	api.sendFoodMessage(chatID, returnMsg)
}

func (api Api) sendFoodMessage(chatID int64, text string) {
	for _, formatter := range api.formatter {
		tmp, err := formatter.Format(text)
		if err != nil {
			text = tmp
		}
	}

	api.tapi.SendMessage(chatID, text)
}
