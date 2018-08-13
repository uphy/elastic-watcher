package trigger

import (
	"errors"
	"fmt"
	"strings"

	"github.com/uphy/elastic-watcher/watcher/jsonutil"
)

type Hourly struct {
	Minute jsonutil.IntArray `json:"minute"`
}

func (c *Hourly) Cron() ([]Cron, error) {
	switch c.Minute.Size() {
	case 0:
		return nil, errors.New("specify `minute` at `hour`")
	case 1:
		return []Cron{newCron("0", fmt.Sprintf("*/%d", c.Minute.Value()[0]), "*", "*", "*", "*")}, nil
	default:
		return []Cron{newCron("0", strings.Join(c.Minute.Strings(), ","), "*", "*", "*", "*")}, nil
	}
}
