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
	wire.Build(config.NewDatabaseConnection)
	return nil, nil
}

func InitTelegramApi() (*telegrambotapi.TelegramBotApi, error) {
	wire.Build(config.NewTelegramBotApi)
	return nil, nil
}

// services

var serviceSet = wire.NewSet(
	InitDatabaseConnection,
	InitTelegramApi,
	db.NewSubscriberRepository,
	domain.NewCliService,
	domain.NewMealService,
	domain.NewMsgService,
	domain.NewSubscriberScheduler,
	domain.NewSubscriberService,
)

func InitMsgService() (domain.MsgService, error) {
	wire.Build(serviceSet)
	return domain.MsgService{}, nil
}
