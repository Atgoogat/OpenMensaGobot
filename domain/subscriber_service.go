package domain

import (
	"errors"
	"log"

	"github.com/Atgoogat/openmensarobot/db"
	"github.com/Atgoogat/openmensarobot/openmensa"
)

var (
	ErrSubscriptionCreation = errors.New("failed to create subscription")
	ErrScheduleSubscription = errors.New("could not schedule subscription")
	ErrUnsubscribeFailed    = errors.New("unsubscribe failed")

	ErrSubscriptionNotFound    = errors.New("subscription not found")
	ErrSubscriptionFetchFailed = errors.New("subscription fetch failed")
)

type SubscriberService struct {
	repo      db.SubscriberRepository
	scheduler *SubscriberScheduler
}

func NewSubscriberService(repo db.SubscriberRepository, scheduler *SubscriberScheduler) SubscriberService {
	return SubscriberService{
		repo:      repo,
		scheduler: scheduler,
	}
}

func (ss SubscriberService) Subscribe(chatID int64, mensaID int, pushHours, pushMinutes uint, priceType openmensa.PriceType) (db.Subscriber, error) {
	log.Printf("update subscriber %d %d", chatID, mensaID)
	sub, err := ss.repo.UpdateSubscriber(chatID, mensaID, pushHours, pushMinutes, priceType)
	if err != nil {
		return sub, errors.Join(ErrSubscriptionCreation, err)
	}

	err = ss.scheduler.InsertJob(sub)
	if err != nil {
		return sub, errors.Join(ErrScheduleSubscription, err)
	}
	return sub, err
}

func (ss SubscriberService) GetSubscription(id uint) (db.Subscriber, error) {
	sub, found, err := ss.repo.FindSubscriberById(id)
	if err != nil {
		return sub, errors.Join(ErrSubscriptionFetchFailed, err)
	}
	if !found {
		return sub, ErrSubscriptionNotFound
	}
	return sub, nil
}

func (ss SubscriberService) Unsubscribe(chatID int64, mensaID int) error {
	err := ss.repo.DeleteSubscriberByChatIDAndMensaID(chatID, mensaID)
	if err != nil {
		return errors.Join(ErrUnsubscribeFailed)
	}
	return nil
}

func (ss SubscriberService) UnsubscribeChat(chatID int64) error {
	err := ss.repo.DeleteSubscriberByChatID(chatID)
	if err != nil {
		return errors.Join(ErrUnsubscribeFailed)
	}
	return nil
}
