package domain

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Atgoogat/openmensarobot/openmensa"
	"github.com/Atgoogat/openmensarobot/telegrambotapi"
)

var (
	ErrInputTooLong      = errors.New("input text is too long")
	ErrCommandNotDefined = errors.New("this command is not defined")

	ErrSubUnexpectedFormat      = errors.New("unexpected argument format. Example: /sub 1719 09:00 students")
	ErrUnsubAllUnexpectedFormat = errors.New("unexpected argument format. Example: /unsuball")
	ErrTodayUnexpectedFormat    = errors.New("unexpected argument format. Example: /today 1719")
)

const (
	MaxTextInputSize = 100
)

type CliService struct {
	subscriberService SubscriberService
	mealService       MealService
}

func NewCliService(subscriberService SubscriberService, mealService MealService) CliService {
	return CliService{
		subscriberService: subscriberService,
		mealService:       mealService,
	}
}

func (cs CliService) ParseAndExecuteCommand(msg telegrambotapi.TelegramMessage) (string, error) {
	if len(msg.Text) > MaxTextInputSize {
		return "", ErrInputTooLong
	}

	if msg.Text == "/start" {
		return "Hello and welcome! Try /help ;)", nil
	} else if strings.HasPrefix(msg.Text, "/sub") {
		return cs.parseAndExecuteSub(msg)
	} else if msg.Text == "/unsuball" {
		return cs.parseAndExecuteUnsubAll(msg)
	} else if strings.HasPrefix(msg.Text, "/today") {
		return cs.parseAndExecuteToday(msg)
	}
	return "", ErrCommandNotDefined
}

var subRegex = regexp.MustCompile(`^/sub \d{1,10} (?:0[0-9]|1[0-9]|2[0-3]):(?:0[0-9]|[1-5][0-9]) (none|students|employees|pupils)$`)

// Parse format: /sub [mensaID] 09:00 priceType
func (cs CliService) parseAndExecuteSub(msg telegrambotapi.TelegramMessage) (string, error) {
	if !subRegex.MatchString(msg.Text) {
		return "", ErrSubUnexpectedFormat
	}

	var mensaID int
	var pushHours, pushMinutes uint
	var price openmensa.PriceType

	n, err := fmt.Sscanf(msg.Text, "/sub %d %d:%d %s", &mensaID, &pushHours, &pushMinutes, &price)
	if n != 4 || err != nil {
		return "", errors.Join(ErrSubUnexpectedFormat, err)
	}

	sub, err := cs.subscriberService.Subscribe(msg.ChatID, mensaID, pushHours, pushMinutes, price)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("subscribed to mensa %d. next update on %02d:%02d", sub.MensaID, sub.Push.Hours, sub.Push.Minutes), nil
}

func (cs CliService) parseAndExecuteUnsubAll(msg telegrambotapi.TelegramMessage) (string, error) {
	err := cs.subscriberService.UnsubscribeChat(msg.ChatID)
	return "everything unsubscribed", err
}

var todayRegex = regexp.MustCompile(`^/today \d{1,10}$`)

// Parse format: /today [mensaID]
func (cs CliService) parseAndExecuteToday(msg telegrambotapi.TelegramMessage) (string, error) {
	if !todayRegex.MatchString(msg.Text) {
		return "", ErrTodayUnexpectedFormat
	}

	var mensaID int
	n, err := fmt.Sscanf(msg.Text, "/today %d", &mensaID)
	if n != 1 || err != nil {
		return "", errors.Join(ErrTodayUnexpectedFormat, err)
	}

	log.Printf("requested mensa today %d %d", msg.ChatID, mensaID)
	return cs.mealService.GetFormatedMeals(mensaID, time.Now(), openmensa.PRICE_NONE)
}
