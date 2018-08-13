package jsonutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type (
	Object struct {
		value interface{}
	}
	Array []Object
)

func (s Object) String() string {
	return fmt.Sprint(s.value)
}

func (s Object) Int() (int, error) {
	return strconv.Atoi(s.String())
}

func (s Object) Value() interface{} {
	return s.value
}

func (o Object) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.value)
}

func (o *Object) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*o = Object{
		value: v,
	}
	return nil
}

func NewArray(v ...interface{}) Array {
	array := Array{}
	for _, e := range v {
		array = append(array, Object{e})
	}
	return array
}

func (a Array) Ints() ([]int, error) {
	s := []int{}
	for _, o := range a {
		i, err := o.Int()
		if err != nil {
			return nil, err
		}
		s = append(s, i)
	}
	return s, nil
}

func (a Array) Strings() []string {
	s := []string{}
	for _, o := range a {
		i := o.String()
		s = append(s, i)
	}
	return s
}

func (a Array) Size() int {
	return len(a)
}

func (a Array) MarshalJSON() ([]byte, error) {
	if a.Size() == 1 {
		return a[0].MarshalJSON()
	}
	v := []Object(a)
	return json.Marshal(v)
}

func (a *Array) UnmarshalJSON(data []byte) error {
	var oo []Object
	if err := json.Unmarshal(data, &oo); err == nil {
		*a = oo
		return nil
	}
	var o Object
	if err := json.Unmarshal(data, &o); err == nil {
		*a = Array{o}
		return nil
	}
	return errors.New("unsupported format array: " + string(data))
}
