package context

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type (
	TemplateValue       string
	TemplateValues      []templateValuesEntry
	templateValuesEntry struct {
		key   string
		value TemplateValue
	}
)

func (t TemplateValue) String(ctx ExecutionContext) (string, error) {
	return renderTemplate(ctx, string(t))
}

func (t TemplateValues) String(ctx ExecutionContext, index int) (string, error) {
	return t[index].value.String(ctx)
}

func (t TemplateValues) StringByKey(ctx ExecutionContext, key string) (string, error) {
	for _, e := range t {
		if e.key == key {
			return e.value.String(ctx)
		}
	}
	return "", nil
}

func (t TemplateValues) Keys() []string {
	keys := []string{}
	for _, e := range t {
		keys = append(keys, e.key)
	}
	return keys
}

func (t TemplateValues) Size() int {
	return len(t)
}

func (t TemplateValues) Map(ctx ExecutionContext) (map[string]string, error) {
	m := map[string]string{}
	for _, key := range t.Keys() {
		v, err := t.StringByKey(ctx, key)
		if err != nil {
			return nil, err
		}
		m[key] = v
	}
	return m, nil
}

func (t TemplateValues) Slice(ctx ExecutionContext) ([]string, error) {
	s := []string{}
	for _, e := range t {
		v, err := e.value.String(ctx)
		if err != nil {
			return nil, err
		}
		s = append(s, v)
	}
	return s, nil
}

func (t TemplateValues) MarshalJSON() ([]byte, error) {
	var v interface{}
	if t.isArray() {
		array := []TemplateValue{}
		for _, e := range t {
			array = append(array, e.value)
		}
		v = array
	} else {
		m := map[string]TemplateValue{}
		for _, e := range t {
			m[e.key] = e.value
		}
		v = m
	}
	return json.Marshal(v)
}

func (t TemplateValues) isArray() bool {
	for _, key := range t.Keys() {
		_, err := strconv.Atoi(key)
		if err != nil {
			return false
		}
	}
	return true
}

func (t *TemplateValues) UnmarshalJSON(data []byte) error {
	// array
	var array []TemplateValue
	if err := json.Unmarshal(data, &array); err == nil {
		tt := TemplateValues{}
		for i, e := range array {
			tt = append(tt, templateValuesEntry{
				key:   fmt.Sprint(i),
				value: e,
			})
		}
		*t = tt
		return nil
	}
	// map
	var m map[string]TemplateValue
	if err := json.Unmarshal(data, &m); err == nil {
		tt := TemplateValues{}
		for i, e := range m {
			tt = append(tt, templateValuesEntry{
				key:   i,
				value: e,
			})
		}
		*t = tt
		return nil
	}
	// single value
	var v TemplateValue
	if err := json.Unmarshal(data, &v); err == nil {
		tt := TemplateValues{}
		tt = append(tt, templateValuesEntry{
			key:   "0",
			value: v,
		})
		*t = tt
		return nil
	}
	return errors.New("unsupported format on TemplateValues:" + string(data))
}
