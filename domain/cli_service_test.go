package domain

import (
	"testing"

	"github.com/Atgoogat/openmensarobot/db"
	"github.com/Atgoogat/openmensarobot/openmensa"
	"github.com/Atgoogat/openmensarobot/telegrambotapi"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func newTestCliService(t *testing.T) (CliService, *SubscriberScheduler) {
	gormDb, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.Migrate(gormDb))
	scheduler := NewSubscriberScheduler()
	subService := NewSubscriberService(db.NewSubscriberRepository(gormDb), scheduler)
	return NewCliService(subService, NewMealService(openmensa.NewOpenmensaApi("http://invalid"))), scheduler
}

func TestUnsub(t *testing.T) {
	cs, scheduler := newTestCliService(t)
	cmd := func(text string) (string, error) {
		return cs.ParseAndExecuteCommand(telegrambotapi.TelegramMessage{ChatID: 42, Text: text})
	}

	_, err := cmd("/sub 1719 09:00 students")
	require.NoError(t, err)
	_, err = cmd("/sub 1720 11:30 none")
	require.NoError(t, err)
	assert.Len(t, scheduler.scheduler.Jobs(), 2)

	res, err := cmd("/unsub 1719")
	assert.NoError(t, err)
	assert.Equal(t, "unsubscribed from mensa 1719", res)
	assert.Len(t, scheduler.scheduler.Jobs(), 1)

	res, err = cmd("/unsub 1719")
	assert.NoError(t, err)
	assert.Equal(t, "you are not subscribed to mensa 1719", res)
	assert.Len(t, scheduler.scheduler.Jobs(), 1)
}

func TestUnsubFormat(t *testing.T) {
	cs, _ := newTestCliService(t)
	for _, text := range []string{"/unsub", "/unsub abc", "/unsub 1 2"} {
		_, err := cs.ParseAndExecuteCommand(telegrambotapi.TelegramMessage{ChatID: 1, Text: text})
		assert.ErrorIs(t, err, ErrUnsubUnexpectedFormat, text)
	}
	// /unsuball must still be routed to unsuball
	res, err := cs.ParseAndExecuteCommand(telegrambotapi.TelegramMessage{ChatID: 1, Text: "/unsuball"})
	assert.NoError(t, err)
	assert.Equal(t, "everything unsubscribed", res)
}
