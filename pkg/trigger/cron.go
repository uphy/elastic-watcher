package trigger

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Cron struct {
	Second    *string `json:"second"`
	Minute    *string `json:"minute"`
	Hour      *string `json:"hour"`
	Day       *string `json:"day"`
	Month     *string `json:"month"`
	DayOfWeek *string `json:"dayofweek"`
	Special   *string `json:"-"`
}

func newCronSpecial(special string) Cron {
	return Cron{
		Special: &special,
	}
}

func newCron(second, minute, hour, day, month, dayOfWeek string) Cron {
	return Cron{&second, &minute, &hour, &day, &month, &dayOfWeek, nil}
}

func (c Cron) Cron() ([]Cron, error) {
	return []Cron{c}, nil
}

func (c Cron) String() string {
	if c.Special != nil {
		return *c.Special
	}
	return fmt.Sprintf("%s %s %s %s %s %s", *c.Second, *c.Minute, *c.Hour, *c.Day, *c.Month, *c.DayOfWeek)
}

func (c Cron) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *Cron) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		ss := strings.Split(s, " ")
		if len(ss) == 6 {
			*c = newCron(ss[0], ss[1], ss[2], ss[3], ss[4], ss[5])
		} else {
			*c = newCronSpecial(s)
		}
	}
	type T Cron
	var t T
	if err := json.Unmarshal(data, &t); err == nil {
		*c = Cron(t)
		return nil
	}
	return errors.New("unsupported cron: " + string(data))
}
