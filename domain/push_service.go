package domain

import (
	"log"
	"time"

	"github.com/Atgoogat/openmensarobot/telegrambotapi"
)

type PushService struct {
	subscriberService SubscriberService
	telegramApi       *telegrambotapi.TelegramBotApi
	mealService       MealService
}

func NewPushService(subscriberService SubscriberService, telegramApi *telegrambotapi.TelegramBotApi, mealService MealService) PushService {
	return PushService{
		subscriberService: subscriberService,
		telegramApi:       telegramApi,
		mealService:       mealService,
	}
}

func (ps PushService) SendPushUpdate(id uint) error {
	log.Printf("push update for %d", id)
	sub, err := ps.subscriberService.GetSubscription(id)
	if err != nil {
		log.Printf("error while getting subscription: %v", err)
		if err == ErrSubscriptionNotFound {
			// return error so that the subscription is removed
			return err
		}
		return nil
	}

	mealsMsg, err := ps.mealService.GetFormatedMeals(sub.MensaID, time.Now(), sub.PriceType)
	if err != nil {
		log.Printf("error while pushing mensa update: %v", err)
		_ = ps.telegramApi.SendHtmlMessage(sub.ChatID, "error while getting todays plan")
		return nil
	}

	_ = ps.telegramApi.SendHtmlMessage(sub.ChatID, mealsMsg)
	return nil
}
