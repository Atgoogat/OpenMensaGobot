package config

import (
	"sync"

	"github.com/Atgoogat/openmensarobot/domain"
)

var (
	scheduler     *domain.SubscriberScheduler
	schedulerOnce sync.Once
)

func GetSubscriberScheduler() *domain.SubscriberScheduler {
	schedulerOnce.Do(func() {
		scheduler = domain.NewSubscriberScheduler()
	})
	return scheduler
}
