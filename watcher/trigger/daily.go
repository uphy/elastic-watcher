package trigger

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/uphy/elastic-watcher/watcher/jsonutil"
)

type (
	Daily struct {
		At Times `json:"at,omitempty"`
	}
	Times []Time
	Time  struct {
		Hour   jsonutil.IntArray `json:"hour"`
		Minute jsonutil.IntArray `json:"minute"`
	}
)

func (c *Daily) Cron() ([]Cron, error) {
	return c.At.Cron()
}

func (t *Time) Cron() ([]Cron, error) {
	return []Cron{newCron("0", strings.Join(t.Minute.Strings(), ","), strings.Join(t.Hour.Strings(), ","), "*", "*", "*")}, nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	if t.Minute.Size() == 1 && t.Hour.Size() == 1 {
		hour := t.Hour.Value()[0]
		minute := t.Minute.Value()[0]
		if minute == 0 {
			if hour == 0 {
				return json.Marshal("midnight")
			} else if hour == 12 {
				return json.Marshal("noon")
			}
		}
		return json.Marshal(fmt.Sprintf("%02d:%02d", hour, minute))
	}
	type T Time
	tt := T(t)
	return json.Marshal(tt)
}

func (t *Time) UnmarshalJSON(data []byte) error {
	/*
		{
		  "trigger" : {
		    "schedule" : {
		      "daily" : {
		        "at" {
		          "hour" : [ 0, 12, 17 ],
		          "minute" : [0, 30]
		        }
		      }
		    }
		  }
		}
	*/
	type T Time
	var tt T
	if err := json.Unmarshal(data, &tt); err == nil {
		*t = Time(tt)
		return nil
	}

	/*
		{
		"trigger" : {
			"schedule" : {
				"daily" : { "at" : "17:00" }
				}
			}
		}
	*/
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		hour, minute, err := parseTime(s)
		if err != nil {
			return err
		}
		*t = Time{
			Hour:   jsonutil.NewIntArray(hour),
			Minute: jsonutil.NewIntArray(minute),
		}
		return nil
	}

	return errors.New("unsupported time format: " + string(data))
}

func (t *Times) Cron() ([]Cron, error) {
	s := []Cron{}
	for _, time := range *t {
		c, err := time.Cron()
		if err != nil {
			return nil, err
		}
		s = append(s, c...)
	}
	return s, nil
}

func (t Times) MarshalJSON() ([]byte, error) {
	if len(t) == 1 {
		return json.Marshal(t[0])
	}
	v := []Time(t)
	return json.Marshal(v)
}

func (t *Times) UnmarshalJSON(data []byte) error {
	var multi []Time
	if err := json.Unmarshal(data, &multi); err == nil {
		*t = Times(multi)
		return nil
	}
	var single Time
	if err := json.Unmarshal(data, &single); err == nil {
		*t = Times{single}
		return nil
	}
	return errors.New("unsupported format times: " + string(data))
}

func parseTime(s string) (hour int, minute int, err error) {
	switch s {
	case "midnight":
		hour = 0
		minute = 0
	case "noon":
		hour = 12
		minute = 0
	default:
		hourMinute := regexp.MustCompile(`^(\d{2}):(\d{2})$`).FindStringSubmatch(s)
		switch len(hourMinute) {
		case 3:
			hour, _ = parseDigit(hourMinute[1])
			minute, _ = parseDigit(hourMinute[2])
		default:
			err = errors.New("invalid time format: " + s)
		}
	}
	return
}

func parseDigit(d string) (int, error) {
	if strings.HasPrefix(d, "0") {
		d = d[1:]
	}
	return strconv.Atoi(d)
}
