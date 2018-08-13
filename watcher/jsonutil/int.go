package jsonutil

import "encoding/json"

type IntArray struct {
	Array
}

func NewIntArray(ints ...int) IntArray {
	v := []interface{}{}
	for _, n := range ints {
		v = append(v, n)
	}
	return IntArray{NewArray(v...)}
}

func (i *IntArray) UnmarshalJSON(data []byte) error {
	var a Array
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	_, err := a.Ints()
	if err != nil {
		return err
	}
	*i = IntArray{a}
	return nil
}

func (i *IntArray) Value() []int {
	ints, err := i.Ints()
	if err != nil {
		panic(err)
	}
	return ints
}
