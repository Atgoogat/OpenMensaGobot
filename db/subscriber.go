package db

import (
	"gorm.io/gorm"
)

type Subscriber struct {
	gorm.Model
	ChatID  int `gorm:"<-:create"`
	MensaID int
	Push    PushTime `gorm:"embedded"`
}

type PushTime struct {
	Hours   uint `gorm:"check:hours < 24"`   // 0 - 23
	Minutes uint `gorm:"check:minutes < 60"` // 0 - 59
}

type SubscriberRepository struct {
	db *gorm.DB
}

func NewSubscriberRepository(db *gorm.DB) SubscriberRepository {
	return SubscriberRepository{
		db: db,
	}
}

func (repo SubscriberRepository) CreateSubscriber(chatID, mensaID int, pushHour, pushMinutes uint) (Subscriber, error) {
	sub := Subscriber{ChatID: chatID, MensaID: mensaID, Push: PushTime{
		Hours:   pushHour,
		Minutes: pushMinutes,
	}}
	res := repo.db.Create(&sub)

	return sub, res.Error
}

func (repo SubscriberRepository) FindSubscriberById(id uint) (Subscriber, error) {
	var sub Subscriber
	result := repo.db.First(&sub, id)
	return sub, result.Error
}

func (repo SubscriberRepository) FindSubscriberByChatID(chatId int) (Subscriber, error) {
	var sub Subscriber
	result := repo.db.First(&sub, "chat_id = ?", chatId)
	return sub, result.Error
}

func (repo SubscriberRepository) DeleteSubscriberByID(id uint) error {
	res := repo.db.Delete(&Subscriber{}, id)
	return res.Error
}
