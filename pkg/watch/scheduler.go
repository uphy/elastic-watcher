package watch

import (
	"github.com/robfig/cron"
)

type (
	Task interface {
		Run()
	}
	Scheduler struct {
		c *cron.Cron
	}
)

func newScheduler() *Scheduler {
	schedule := &Scheduler{
		c: cron.New(),
	}
	return schedule
}

func (s *Scheduler) AddTask(schedule string, task Task) error {
	return s.c.AddJob(schedule, task)
}

func (s *Scheduler) Start() {
	s.c.Start()
}

func (s *Scheduler) Stop() {
	s.c.Stop()
}
