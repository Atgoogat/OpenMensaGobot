package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var repo SubscriberRepository

func TestMain(m *testing.M) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = Migrate(db)
	if err != nil {
		panic(err)
	}
	repo = NewSubscriberRepository(db)
	m.Run()
}

func TestCreateSubscriber(t *testing.T) {
	sub, err := repo.CreateSubscriber(1, 2, 11, 30)
	assert.Nil(t, err)
	sub2, err := repo.CreateSubscriber(2, 3, 11, 31)
	assert.Nil(t, err)
	assert.NotEqual(t, sub.ID, sub2.ID)
}

func TestCreateAndRetrieveSubscriber(t *testing.T) {
	sub, err := repo.CreateSubscriber(12, 13, 10, 00)
	assert.Nil(t, err)

	id := sub.ID
	sub, err = repo.FindSubscriberById(id)
	assert.Nil(t, err)

	assert.Equal(t, id, sub.ID)
	assert.Equal(t, 12, sub.ChatID)
	assert.Equal(t, 13, sub.MensaID)
	assert.Equal(t, uint(10), sub.Push.Hours)
	assert.Equal(t, uint(0), sub.Push.Minutes)
}

func TestInsertWrongTime(t *testing.T) {
	_, err := repo.CreateSubscriber(0, 0, 24, 0)
	assert.NotNil(t, err)
	_, err = repo.CreateSubscriber(0, 0, 0, 60)
	assert.NotNil(t, err)
}

func TestCreateAndRetrieveByChatID(t *testing.T) {
	_, err := repo.CreateSubscriber(123, 0, 0, 0)
	assert.Nil(t, err)

	sub, err := repo.FindSubscriberByChatID(123)
	assert.Nil(t, err)
	assert.Equal(t, 123, sub.ChatID)
	assert.Equal(t, 0, sub.MensaID)
}
