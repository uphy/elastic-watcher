package trigger

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/uphy/elastic-watcher/pkg/jsonutil"
)

func TestCron(t *testing.T) {
	tests := []struct {
		name    string
		cronner Cronner
		expect  []Cron
		err     bool
	}{
		{
			name: "hourly",
			cronner: &Hourly{
				Minute: jsonutil.NewIntArray(10),
			},
			expect: []Cron{newCron("0", "*/10", "*", "*", "*", "*")},
		},
		{
			name: "hourly multi",
			cronner: &Hourly{
				Minute: jsonutil.NewIntArray(10, 20),
			},
			expect: []Cron{newCron("0", "10,20", "*", "*", "*", "*")},
		},
		{
			name: "hourly empty",
			cronner: &Hourly{
				Minute: jsonutil.NewIntArray(),
			},
			err: true,
		},
		{
			name: "time",
			cronner: &Time{
				Hour:   jsonutil.NewIntArray(1),
				Minute: jsonutil.NewIntArray(10),
			},
			expect: []Cron{newCron("0", "10", "1", "*", "*", "*")},
		},
		{
			name: "time multi",
			cronner: &Time{
				Hour:   jsonutil.NewIntArray(0, 12),
				Minute: jsonutil.NewIntArray(0, 30),
			},
			expect: []Cron{
				newCron("0", "0,30", "0,12", "*", "*", "*"),
			},
		},
		{
			name: "times",
			cronner: &Times{
				Time{
					Hour:   jsonutil.NewIntArray(0, 12),
					Minute: jsonutil.NewIntArray(0, 30),
				},
				Time{
					Hour:   jsonutil.NewIntArray(17),
					Minute: jsonutil.NewIntArray(0),
				},
			},
			expect: []Cron{
				newCron("0", "0,30", "0,12", "*", "*", "*"),
				newCron("0", "0", "17", "*", "*", "*"),
			},
		},
		{
			name: "daily",
			cronner: &Daily{
				At: Times{
					Time{
						Hour:   jsonutil.NewIntArray(0, 12),
						Minute: jsonutil.NewIntArray(0, 30),
					},
					Time{
						Hour:   jsonutil.NewIntArray(17),
						Minute: jsonutil.NewIntArray(0),
					},
				},
			},
			expect: []Cron{
				newCron("0", "0,30", "0,12", "*", "*", "*"),
				newCron("0", "0", "17", "*", "*", "*"),
			},
		},
		{
			name: "weekly",
			cronner: &Weekly{
				Daily: Daily{
					At: Times{
						Time{
							Hour:   jsonutil.NewIntArray(0, 12),
							Minute: jsonutil.NewIntArray(0, 30),
						},
						Time{
							Hour:   jsonutil.NewIntArray(17),
							Minute: jsonutil.NewIntArray(0),
						},
					},
				},
				On: jsonutil.NewStringArray("Sunday", "Saturday"),
			},
			expect: []Cron{
				newCron("0", "0,30", "0,12", "*", "*", "sun,sat"),
				newCron("0", "0", "17", "*", "*", "sun,sat"),
			},
		},
		{
			name: "monthly",
			cronner: &Monthly{
				Daily: Daily{
					At: Times{
						Time{
							Hour:   jsonutil.NewIntArray(0, 12),
							Minute: jsonutil.NewIntArray(0, 30),
						},
						Time{
							Hour:   jsonutil.NewIntArray(17),
							Minute: jsonutil.NewIntArray(0),
						},
					},
				},
				On: jsonutil.NewIntArray(1, 15),
			},
			expect: []Cron{
				newCron("0", "0,30", "0,12", "1,15", "*", "*"),
				newCron("0", "0", "17", "1,15", "*", "*"),
			},
		},
		{
			name: "yearly",
			cronner: &Yearly{
				Monthly: Monthly{
					Daily: Daily{
						At: Times{
							Time{
								Hour:   jsonutil.NewIntArray(0, 12),
								Minute: jsonutil.NewIntArray(0, 30),
							},
							Time{
								Hour:   jsonutil.NewIntArray(17),
								Minute: jsonutil.NewIntArray(0),
							},
						},
					},
					On: jsonutil.NewIntArray(1, 15),
				},
				In: jsonutil.NewStringArray("january", "3"),
			},
			expect: []Cron{
				newCron("0", "0,30", "0,12", "1,15", "1,3", "*"),
				newCron("0", "0", "17", "1,15", "1,3", "*"),
			},
		},
	}
	for _, test := range tests {
		c, err := test.cronner.Cron()
		if test.err {
			if err == nil {
				t.Errorf("expect err: %s", test.name)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("failed on Cron(); %s : %s", err, test.name)
			}
			if len(test.expect) != len(c) {
				t.Errorf("expect %d but %d: %s", len(test.expect), len(c), test.name)
				continue
			}
			for i, e := range test.expect {
				if !reflect.DeepEqual(e, c[i]) {
					t.Errorf("expect %s but %s: %s", e, c[i], test.name)
				}
			}
		}
	}
}

func TestTriggerMarshal(t *testing.T) {
	tests := []struct {
		name   string
		JSON   string
		Expect string
	}{
		{
			name: "monthly single",
			JSON: `{
					"schedule": {
						"monthly":{"at":["midnight","noon"],"on":[10,20]}
					}
				}`,
		},
		{
			name: "monthly array",
			JSON: `{
					"schedule": {
						"monthly":[
							{"at":["midnight","noon"],"on":[10,20]},
							{"at":["13:00","07:00"],"on":[1,25]}
						]
					}
				}`,
		},
		{
			name: "yearly single",
			JSON: `{
					"schedule": {
						"yearly":{"in":"janualy","at":["midnight","noon"],"on":[10,20]}
					}
				}`,
		},
		{
			name: "yearly array",
			JSON: `{
					"schedule": {
						"yearly":[
							{"in":"janualy" ,"at":["midnight","noon"],"on":[10,20]},
							{"in":"janualy" ,"at":["13:00","07:00"],"on":[1,25]}
						]
					}
				}`,
		},
	}
	for _, test := range tests {
		var result Trigger
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

		var format = func(s string) string {
			var m map[string]interface{}
			if err := json.Unmarshal([]byte(s), &m); err != nil {
				t.Errorf("failed to unmarshal on `%s`: %s", test.name, s)
				return s
			}
			formatted, err := json.MarshalIndent(m, "", "   ")
			if err != nil {
				t.Errorf("failed to marshal on `%s`: %s", test.name, s)
				return s
			}
			return string(formatted)
		}
		expectFormatted := format(expect)
		actualFormatted := format(actual)
		if !reflect.DeepEqual(expectFormatted, actualFormatted) {
			fmt.Println(expectFormatted)
			fmt.Println(actualFormatted)
			t.Errorf("expect %#v but %#v: %s", expect, actual, test.name)
		}
	}
}
