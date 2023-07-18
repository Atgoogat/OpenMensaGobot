package db

import (
	"github.com/Atgoogat/openmensarobot/openmensa"
	"gorm.io/gorm"
)

type Subscriber struct {
	gorm.Model
	ChatID    int64 `gorm:"<-:create"`
	MensaID   int
	Push      PushTime `gorm:"embedded"`
	PriceType openmensa.PriceType
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

func (repo SubscriberRepository) CreateSubscriber(chatID int64, mensaID int, pushHour, pushMinutes uint, priceType openmensa.PriceType) (Subscriber, error) {
	sub := Subscriber{ChatID: chatID, MensaID: mensaID, Push: PushTime{
		Hours:   pushHour,
		Minutes: pushMinutes,
	}, PriceType: priceType}
	res := repo.db.Create(&sub)

	return sub, res.Error
}

func (repo SubscriberRepository) UpdateSubscriber(chatID int64, mensaID int, pushHour, pushMinutes uint, priceType openmensa.PriceType) (Subscriber, error) {
	sub, found, err := repo.FindSubscriberByChatIDAndMensaID(chatID, mensaID)
	if err != nil {
		return sub, err
	}

	if found {
		sub.Push.Hours = pushHour
		sub.Push.Minutes = pushMinutes
		sub.PriceType = priceType
	} else {
		sub = Subscriber{ChatID: chatID, MensaID: mensaID, Push: PushTime{
			Hours:   pushHour,
			Minutes: pushMinutes,
		}, PriceType: priceType}
	}

	res := repo.db.Save(&sub)

	return sub, res.Error
}

func firstOrDefault[T any](input []T) (ret T) {
	if len(input) > 0 {
		ret = input[0]
	}
	return
}

func (repo SubscriberRepository) FindSubscriberById(id uint) (Subscriber, bool, error) {
	var sub []Subscriber
	result := repo.db.Find(&sub, id).Limit(1)
	return firstOrDefault(sub), result.RowsAffected > 0, result.Error
}

func (repo SubscriberRepository) FindSubscriberByChatID(chatID int64) (Subscriber, error) {
	var sub Subscriber
	result := repo.db.First(&sub, "chat_id = ?", chatID)
	return sub, result.Error
}

func (repo SubscriberRepository) FindSubscriberByChatIDAndMensaID(chatID int64, mensaID int) (Subscriber, bool, error) {
	var sub []Subscriber
	result := repo.db.Find(&sub, "chat_id = ? AND mensa_id = ?", chatID, mensaID).Limit(1)
	return firstOrDefault(sub), result.RowsAffected > 0, result.Error
}

func (repo SubscriberRepository) ExistsSubscriberWithChatIDAndMensaID(chatID int64, mensaID int) (bool, error) {
	var sub []Subscriber
	result := repo.db.Find(&sub, "chat_id = ? AND mensa_id = ?", chatID, mensaID).Limit(1)
	return result.RowsAffected > 0, result.Error
}

func (repo SubscriberRepository) DeleteSubscriberByID(id uint) error {
	res := repo.db.Delete(&Subscriber{}, id)
	return res.Error
}

func (repo SubscriberRepository) DeleteSubscriberByChatIDAndMensaID(chatID int64, mensaID int) error {
	res := repo.db.Where("chat_id = ? AND mensa_id = ?", chatID, mensaID).Delete(&Subscriber{})
	return res.Error
}

func (repo SubscriberRepository) DeleteSubscriberByChatID(chatID int64) error {
	res := repo.db.Where("chat_id = ?", chatID).Delete(&Subscriber{})
	return res.Error
}

func (repo SubscriberRepository) FindAllSubscriber() ([]Subscriber, error) {
	var subs []Subscriber
	res := repo.db.Find(&subs)
	return subs, res.Error
}
