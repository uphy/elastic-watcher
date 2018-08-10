package watcher

import (
	"errors"

	"github.com/robfig/cron"
)

type (
	Task interface {
		Run()
	}
	Scheduler struct {
		c *cron.Cron
	}
	Schedule struct {
		Cron     *string    `json:"cron,omitempty"`
		Interval *string    `json:"interval,omitempty"`
		Hourly   *Hourly    `json:"hourly,omitempty"`
		Daily    *Daily     `json:"daily,omitempty"`
		Weekly   *[]Weekly  `json:"weekly,omitempty"`
		Monthly  *[]Monthly `json:"monthly,omitempty"`
		Yearly   *[]Yearly  `json:"yearly,omitempty"`
	}
	Hourly struct {
		Minute []string `json:"minute"`
	}
	Daily struct {
		At interface{} `json:"at,omitempty"`
	}
	Weekly struct {
		On []string `json:"on"`
		At []string `json:"at"`
	}
	Monthly struct {
		On []string `json:"on"`
		At []string `json:"at"`
	}
	Yearly struct {
		In []string `json:"in"`
		On []string `json:"on"`
		At []string `json:"at"`
	}
)

func newScheduler() *Scheduler {
	schedule := &Scheduler{
		c: cron.New(),
	}
	return schedule
}

func (s *Scheduler) AddTask(schedule *Schedule, task Task) error {
	if schedule.Cron != nil {
		s.c.AddJob(*schedule.Cron, task)
		return nil
	}
	return errors.New("unsupported schedule type")
}

func (s *Scheduler) Start() {
	s.c.Start()
}

func (s *Scheduler) Stop() {
	s.c.Stop()
}
