package trigger

import (
	"encoding/json"
	"testing"
)

func TestParseTime(t *testing.T) {
	tests := []struct {
		time   string
		hour   int
		minute int
		err    bool
	}{
		{
			time:   "00:12",
			hour:   0,
			minute: 12,
		},
		{
			time:   "noon",
			hour:   12,
			minute: 0,
		},
		{
			time:   "midnight",
			hour:   0,
			minute: 0,
		},
		{
			time: "0:0",
			err:  true,
		},
		{
			time: "aaa",
			err:  true,
		},
	}
	for _, test := range tests {
		hour, minute, err := parseTime(test.time)
		if test.err && err == nil {
			t.Error("expect error: " + test.time)
		} else {
			if hour != test.hour || minute != test.minute {
				t.Errorf("expect %2d:%2d but %2d:%2d", test.hour, test.minute, hour, minute)
			}
		}
	}
}

func TestMarshalTimes(t *testing.T) {
	tests := []struct {
		name   string
		JSON   string
		Expect string
	}{
		{
			name: "time",
			JSON: `"17:00"`,
		},
		{
			name: "times",
			JSON: `["midnight","noon","17:00"]`,
		},
		{
			name:   "time json",
			JSON:   `{"hour":17,"minute":0}`,
			Expect: `"17:00"`,
		},
		{
			name: "times json",
			JSON: `{"hour":[0,12],"minute":[0,30]}`,
		},
	}
	for _, test := range tests {
		var result Times
		if err := json.Unmarshal([]byte(test.JSON), &result); err != nil {
			t.Error(err)
			continue
		}
		b, err := json.Marshal(result)
		if err != nil {
			t.Error(err)
			continue
		}
		expect := test.Expect
		if expect == "" {
			expect = test.JSON
		}
		actual := string(b)
		if expect != actual {
			t.Errorf("expect %#v but %#v: %s", expect, actual, test.name)
		}
	}
}
