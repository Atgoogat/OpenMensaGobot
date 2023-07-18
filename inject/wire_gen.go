// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package inject

import (
	"github.com/Atgoogat/openmensarobot/config"
	"github.com/Atgoogat/openmensarobot/db"
	"github.com/Atgoogat/openmensarobot/domain"
	"github.com/Atgoogat/openmensarobot/telegrambotapi"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// Injectors from wire.go:

func InitDatabaseConnection() (*gorm.DB, error) {
	db, err := config.GetDatabaseConnection()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitTelegramApi() (*telegrambotapi.TelegramBotApi, error) {
	telegramBotApi, err := config.GetTelegramBotApi()
	if err != nil {
		return nil, err
	}
	return telegramBotApi, nil
}

func InitMsgService() (domain.MsgService, error) {
	telegramBotApi, err := InitTelegramApi()
	if err != nil {
		return domain.MsgService{}, err
	}
	gormDB, err := InitDatabaseConnection()
	if err != nil {
		return domain.MsgService{}, err
	}
	subscriberRepository := db.NewSubscriberRepository(gormDB)
	subscriberScheduler := config.GetSubscriberScheduler()
	subscriberService := domain.NewSubscriberService(subscriberRepository, subscriberScheduler)
	openmensaApi := config.GetOpenmensaApi()
	mealService := config.GetMealService(openmensaApi)
	cliService := domain.NewCliService(subscriberService, mealService)
	msgService := domain.NewMsgService(telegramBotApi, cliService)
	return msgService, nil
}

func InitPushService() (domain.PushService, error) {
	gormDB, err := InitDatabaseConnection()
	if err != nil {
		return domain.PushService{}, err
	}
	subscriberRepository := db.NewSubscriberRepository(gormDB)
	subscriberScheduler := config.GetSubscriberScheduler()
	subscriberService := domain.NewSubscriberService(subscriberRepository, subscriberScheduler)
	telegramBotApi, err := InitTelegramApi()
	if err != nil {
		return domain.PushService{}, err
	}
	openmensaApi := config.GetOpenmensaApi()
	mealService := config.GetMealService(openmensaApi)
	pushService := domain.NewPushService(subscriberService, telegramBotApi, mealService)
	return pushService, nil
}

func InitScheduler() (*domain.SubscriberScheduler, error) {
	subscriberScheduler := config.GetSubscriberScheduler()
	return subscriberScheduler, nil
}

// wire.go:

var serviceSet = wire.NewSet(

	InitDatabaseConnection,
	InitTelegramApi, config.GetOpenmensaApi, config.GetSubscriberScheduler, config.GetMealService, db.NewSubscriberRepository, domain.NewCliService, domain.NewMsgService, domain.NewSubscriberService, domain.NewPushService,
)
