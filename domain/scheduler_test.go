package domain

import (
	"sync"
	"testing"

	"github.com/Atgoogat/openmensarobot/db"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSchedulerConcurrentInsertRemove(t *testing.T) {
	s := NewSubscriberScheduler()
	var wg sync.WaitGroup
	for i := 1; i <= 50; i++ {
		wg.Add(1)
		go func(id uint) {
			defer wg.Done()
			sub := db.Subscriber{Model: gorm.Model{ID: id}, Push: db.PushTime{Hours: 9, Minutes: 0}}
			assert.NoError(t, s.InsertJob(sub))
			assert.NoError(t, s.InsertJob(sub)) // re-insert replaces
			if id%2 == 0 {
				s.RemoveJob(id)
			}
		}(uint(i))
	}
	wg.Wait()
	assert.Len(t, s.scheduledJobs, 25)
	assert.Len(t, s.scheduler.Jobs(), 25)
}
