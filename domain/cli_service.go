package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Atgoogat/openmensarobot/openmensa"
	"github.com/Atgoogat/openmensarobot/telegrambotapi"
)

var (
	ErrInputTooLong      = errors.New("input text is too long")
	ErrCommandNotDefined = errors.New("this command is not defined")

	ErrSubUnexpectedFormat = errors.New("unexpected argument format. Example: /sub 1719 09:00 students")

	ErrUnsubAllUnexpectedFormat = errors.New("unexpected argument format. Example: /unsuball")
)

const (
	MaxTextInputSize = 100
)

type CliService struct {
	ss SubscriberService
}

func NewCliService(ss SubscriberService) CliService {
	return CliService{
		ss: ss,
	}
}

func (cs CliService) ParseAndExecuteCommand(msg telegrambotapi.TelegramMessage) (string, error) {
	if len(msg.Text) > MaxTextInputSize {
		return "", ErrInputTooLong
	}

	if strings.HasPrefix(msg.Text, "/sub") {
		return cs.parseAndExecuteSub(msg)
	} else if msg.Text == "/unsuball" {
		return cs.parseAndExecuteUnsubAll(msg)
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
		return "", ErrSubUnexpectedFormat
	}

	sub, err := cs.ss.Subscribe(msg.ChatID, mensaID, pushHours, pushMinutes, price)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("subscribed to mensa %d. next update on %02d:%02d", sub.MensaID, sub.Push.Hours, sub.Push.Minutes), nil
}

func (cs CliService) parseAndExecuteUnsubAll(msg telegrambotapi.TelegramMessage) (string, error) {
	err := cs.ss.UnsubscribeChat(msg.ChatID)
	return "everything unsubscribed", err
}
