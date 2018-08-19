package trigger

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/uphy/elastic-watcher/pkg/jsonutil"
)

type Monthlys []Monthly

type Monthly struct {
	Daily
	On jsonutil.IntArray `json:"on"`
}

func (c *Monthly) Cron() ([]Cron, error) {
	cron, err := c.At.Cron()
	if err != nil {
		return nil, err
	}
	for i, e := range cron {
		day := strings.Join(c.On.Strings(), ",")
		e.Day = &day
		cron[i] = e
	}
	return cron, nil
}

func (m Monthlys) MarshalJSON() ([]byte, error) {
	if len(m) == 1 {
		return json.Marshal(m[0])
	}
	v := []Monthly(m)
	return json.Marshal(v)
}

func (m *Monthlys) UnmarshalJSON(data []byte) error {
	var array []Monthly
	if err := json.Unmarshal(data, &array); err == nil {
		*m = Monthlys(array)
		return nil
	}
	var single Monthly
	if err := json.Unmarshal(data, &single); err == nil {
		*m = Monthlys{single}
		return nil
	}
	return errors.New("unsupported monthly format: " + string(data))
}
