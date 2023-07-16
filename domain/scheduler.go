package domain

import (
	"fmt"
	"time"

	"github.com/Atgoogat/openmensarobot/db"
	"github.com/go-co-op/gocron"
)

// Runs everytime a subscriber is scheduled according to its pushtime
// parameter is the subscriber id
// if an error is returned the subscriber is removed form the scheduler
type SubscriberScheduleAction func(uint) error

type SubscriberScheduler struct {
	scheduler     *gocron.Scheduler
	scheduledJobs map[uint]*gocron.Job
	action        SubscriberScheduleAction
}

func NewSubscriberScheduler(action SubscriberScheduleAction) *SubscriberScheduler {
	scheduler := gocron.NewScheduler(time.Local)
	scheduler.StartAsync()

	return &SubscriberScheduler{
		scheduler:     scheduler,
		scheduledJobs: make(map[uint]*gocron.Job),
		action:        action,
	}
}

func (scheduler *SubscriberScheduler) InsertJob(subscriber db.Subscriber) error {
	time := fmt.Sprintf("%02d:%02d", subscriber.Push.Hours, subscriber.Push.Minutes)
	id := subscriber.ID
	job, err := scheduler.scheduler.Every(1).Day().At(time).Do(func() {
		if scheduler.action != nil {
			err := scheduler.action(id)
			if err != nil {
				j := scheduler.scheduledJobs[id]
				scheduler.scheduler.Remove(j)
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
	job := scheduler.scheduledJobs[id]
	scheduler.scheduler.Remove(job)
}
