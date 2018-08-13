package trigger

import (
	"errors"
	"strings"

	"github.com/uphy/elastic-watcher/watcher/jsonutil"
)

type Weekly struct {
	Daily
	On jsonutil.StringArray `json:"on"`
}

func (c *Weekly) Cron() ([]Cron, error) {
	ons := []string{}
	for _, s := range c.On.Value() {
		var on string
		switch strings.ToLower(s) {
		case "1", "sun", "sunday":
			on = "sun"
		case "2", "mon", "monday":
			on = "mon"
		case "3", "tue", "tuesday":
			on = "tue"
		case "4", "wed", "wednesday":
			on = "wed"
		case "5", "thu", "thursday":
			on = "thu"
		case "6", "fri", "friday":
			on = "fri"
		case "7", "sat", "saturday":
			on = "sat"
		default:
			return nil, errors.New("unsupported weekly `on` format: " + s)
		}
		ons = append(ons, on)
	}
	at, err := c.At.Cron()
	if err != nil {
		return nil, err
	}
	dayOfWeek := strings.Join(ons, ",")
	for i, a := range at {
		a.DayOfWeek = &dayOfWeek
		at[i] = a
	}
	return at, nil
}
