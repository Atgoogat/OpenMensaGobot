package db

import (
	"gorm.io/gorm"
)

type Subscriber struct {
	gorm.Model
	ChatID  int64 `gorm:"<-:create"`
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

func (repo SubscriberRepository) CreateSubscriber(chatID int64, mensaID int, pushHour, pushMinutes uint) (Subscriber, error) {
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

func (repo SubscriberRepository) FindSubscriberByChatID(chatID int64) (Subscriber, error) {
	var sub Subscriber
	result := repo.db.First(&sub, "chat_id = ?", chatID)
	return sub, result.Error
}

func (repo SubscriberRepository) FindSubscriberByChatIDAndMensaID(chatID int64, mensaID int) (Subscriber, error) {
	var sub Subscriber
	result := repo.db.First(&sub, "chat_id = ? AND mensa_id = ?", chatID, mensaID)
	return sub, result.Error
}

func (repo SubscriberRepository) ExistsSubscriberWithChatIDAndMensaID(chatID int64, mensaID int) (bool, error) {
	var sub []Subscriber
	result := repo.db.Find(&sub, "chat_id = ? AND mensa_id = ?", chatID, mensaID)
	return result.RowsAffected > 0, result.Error
}

func (repo SubscriberRepository) DeleteSubscriberByID(id uint) error {
	res := repo.db.Delete(&Subscriber{}, id)
	return res.Error
}
