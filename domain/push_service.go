package domain

import (
	"errors"
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
		if errors.Is(err, ErrSubscriptionNotFound) {
			// return error so that the subscription is removed
			return err
		}
		return nil
	}

	mealsMsg, err := ps.mealService.GetFormatedMeals(sub.MensaID, time.Now(), sub.PriceType)
	if errors.Is(err, ErrNoMeals) {
		// mensa closed (weekend, holiday) - nothing to push
		log.Printf("no meals for mensa %d today, skipping push for %d", sub.MensaID, id)
		return nil
	}
	if err != nil {
		log.Printf("error while pushing mensa update: %v", err)
		err = ps.telegramApi.SendHtmlMessage(sub.ChatID, "error while getting todays plan")
	} else {
		err = ps.telegramApi.SendHtmlMessage(sub.ChatID, mealsMsg)
	}

	if telegrambotapi.IsChatUnreachable(err) {
		log.Printf("chat %d is unreachable, removing subscription %d", sub.ChatID, id)
		if delErr := ps.subscriberService.repo.DeleteSubscriberByID(id); delErr != nil {
			log.Printf("could not delete subscription %d: %v", id, delErr)
		}
		// return error so that the job is removed
		return err
	}
	return nil
}
