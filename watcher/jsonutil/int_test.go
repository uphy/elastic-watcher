package jsonutil

import (
	"encoding/json"
	"testing"
)

func TestIntArray(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{
			name: "number array",
			json: "[0,1,2]",
		},
	}
	for _, test := range tests {
		var v IntArray
		if err := json.Unmarshal([]byte(test.json), &v); err != nil {
			t.Error(err)
		}
		s, err := json.Marshal(v)
		if err != nil {
			t.Error(err)
		}
		if string(s) != test.json {
			t.Errorf("expect %s but %s", test.json, string(s))
		}
	}
}
