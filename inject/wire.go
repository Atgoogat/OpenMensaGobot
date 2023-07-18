//go:build wireinject
// +build wireinject

package inject

import (
	"github.com/Atgoogat/openmensarobot/config"
	"github.com/Atgoogat/openmensarobot/db"
	"github.com/Atgoogat/openmensarobot/domain"
	"github.com/Atgoogat/openmensarobot/telegrambotapi"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// external resources

func InitDatabaseConnection() (*gorm.DB, error) {
	wire.Build(config.GetDatabaseConnection)
	return nil, nil
}

func InitTelegramApi() (*telegrambotapi.TelegramBotApi, error) {
	wire.Build(config.GetTelegramBotApi)
	return nil, nil
}

// services

var serviceSet = wire.NewSet(
	// singletons
	InitDatabaseConnection,
	InitTelegramApi,
	config.GetOpenmensaApi,
	config.GetSubscriberScheduler,

	config.GetTextFormatter,
	db.NewSubscriberRepository,
	domain.NewMealService,
	domain.NewCliService,
	domain.NewMsgService,
	domain.NewSubscriberService,
	domain.NewPushService,
)

func InitMsgService() (domain.MsgService, error) {
	wire.Build(serviceSet)
	return domain.MsgService{}, nil
}

func InitPushService() (domain.PushService, error) {
	wire.Build(serviceSet)
	return domain.PushService{}, nil
}

func InitScheduler() (*domain.SubscriberScheduler, error) {
	wire.Build(serviceSet)
	return nil, nil
}
