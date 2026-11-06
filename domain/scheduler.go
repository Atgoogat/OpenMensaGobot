package domain

import (
	"fmt"
	"sync"
	"time"

	"github.com/Atgoogat/openmensarobot/db"
	"github.com/go-co-op/gocron"
)

// Runs everytime a subscriber is scheduled according to its pushtime
// parameter is the subscriber id
// if an error is returned the subscriber is removed form the scheduler
type SubscriberScheduleAction func(uint) error

type SubscriberScheduler struct {
	scheduler *gocron.Scheduler
	action    SubscriberScheduleAction

	mu            sync.Mutex
	scheduledJobs map[uint]*gocron.Job
}

func NewSubscriberScheduler() *SubscriberScheduler {
	scheduler := gocron.NewScheduler(time.Local)
	scheduler.StartAsync()

	return &SubscriberScheduler{
		scheduler:     scheduler,
		scheduledJobs: make(map[uint]*gocron.Job),
	}
}

func (scheduler *SubscriberScheduler) SetAction(action SubscriberScheduleAction) {
	scheduler.action = action
}

func (scheduler *SubscriberScheduler) InsertJob(subscriber db.Subscriber) error {
	scheduler.mu.Lock()
	defer scheduler.mu.Unlock()

	scheduler.removeJobLocked(subscriber.ID)

	at := fmt.Sprintf("%02d:%02d", subscriber.Push.Hours, subscriber.Push.Minutes)
	id := subscriber.ID
	var job *gocron.Job
	job, err := scheduler.scheduler.Every(1).Day().At(at).Do(func() {
		if scheduler.action == nil {
			return
		}
		if err := scheduler.action(id); err != nil {
			scheduler.mu.Lock()
			defer scheduler.mu.Unlock()
			// only remove if the job was not replaced in the meantime
			if scheduler.scheduledJobs[id] == job {
				scheduler.removeJobLocked(id)
			}
		}
	})
	if err != nil {
		return err
	}
	scheduler.scheduledJobs[id] = job
	return nil
}

func (scheduler *SubscriberScheduler) RemoveJob(id uint) {
	scheduler.mu.Lock()
	defer scheduler.mu.Unlock()
	scheduler.removeJobLocked(id)
}

func (scheduler *SubscriberScheduler) removeJobLocked(id uint) {
	if job, ok := scheduler.scheduledJobs[id]; ok {
		scheduler.scheduler.RemoveByReference(job)
		delete(scheduler.scheduledJobs, id)
	}
}
