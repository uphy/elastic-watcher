package trigger

import (
	"github.com/pkg/errors"
)

type (
	Trigger struct {
		Schedule *Schedule `json:"schedule"`
	}
	Schedule struct {
		Cron     *Cron     `json:"cron,omitempty"`
		Interval *Interval `json:"interval,omitempty"`
		Hourly   *Hourly   `json:"hourly,omitempty"`
		Daily    *Daily    `json:"daily,omitempty"`
		Weekly   *Weekly   `json:"weekly,omitempty"`
		Monthly  *Monthlys `json:"monthly,omitempty"`
		Yearly   *Yearlys  `json:"yearly,omitempty"`
	}
	Cronner interface {
		Cron() ([]Cron, error)
	}
)

func (s *Schedule) CronSchedules() ([]string, error) {
	// collect cronners
	cronners := []Cronner{}
	if s.Cron != nil {
		cronners = append(cronners, s.Cron)
	}
	if s.Interval != nil {
		cronners = append(cronners, s.Interval)
	}
	if s.Hourly != nil {
		cronners = append(cronners, s.Hourly)
	}
	if s.Daily != nil {
		cronners = append(cronners, s.Daily)
	}
	if s.Weekly != nil {
		cronners = append(cronners, s.Weekly)
	}
	if s.Monthly != nil {
		for _, m := range *s.Monthly {
			cronners = append(cronners, &m)
		}
	}
	if s.Yearly != nil {
		for _, m := range *s.Yearly {
			cronners = append(cronners, &m)
		}
	}
	// collect crons
	crons := []Cron{}
	for _, cron := range cronners {
		if cron != nil {
			c, err := cron.Cron()
			if err != nil {
				return nil, errors.Wrapf(err, "failed to generate cron schedule: %#v", cron)
			}
			crons = append(crons, c...)
		}
	}

	// convert crons to string
	c := []string{}
	for _, e := range crons {
		c = append(c, e.String())
	}
	return c, nil
}
