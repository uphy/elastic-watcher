package trigger

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/uphy/elastic-watcher/pkg/jsonutil"
)

type Yearlys []Yearly
type Yearly struct {
	Monthly
	In jsonutil.StringArray `json:"in"`
}

func (c *Yearly) Cron() ([]Cron, error) {
	cron, err := c.Monthly.Cron()
	if err != nil {
		return nil, err
	}
	months := []string{}
	for _, month := range c.In.Strings() {
		var m string
		switch strings.ToLower(month) {
		case "1", "january", "jan":
			m = "1"
		case "2", "february", "feb":
			m = "2"
		case "3", "march", "mar":
			m = "3"
		case "4", "april", "apr":
			m = "4"
		case "5", "may":
			m = "5"
		case "6", "june", "jun":
			m = "6"
		case "7", "july", "jul":
			m = "7"
		case "8", "august", "aug":
			m = "8"
		case "9", "september", "sep":
			m = "9"
		case "10", "october", "oct":
			m = "10"
		case "11", "november", "nov":
			m = "11"
		case "12", "december", "dec":
			m = "12"
		default:
			return nil, errors.New("unsupported month format: " + month)
		}
		months = append(months, m)
	}
	m := strings.Join(months, ",")
	for i, e := range cron {
		month := m
		e.Month = &month
		cron[i] = e
	}
	return cron, nil
}

func (m Yearlys) MarshalJSON() ([]byte, error) {
	if len(m) == 1 {
		return json.Marshal(m[0])
	}
	v := []Yearly(m)
	return json.Marshal(v)
}

func (m *Yearlys) UnmarshalJSON(data []byte) error {
	var array []Yearly
	if err := json.Unmarshal(data, &array); err == nil {
		*m = Yearlys(array)
		return nil
	}
	var single Yearly
	if err := json.Unmarshal(data, &single); err == nil {
		*m = Yearlys{single}
		return nil
	}
	return errors.New("unsupported yearly format: " + string(data))
}
