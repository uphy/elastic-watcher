package jsonutil

type StringArray struct {
	Array
}

func NewStringArray(s ...string) StringArray {
	v := []interface{}{}
	for _, n := range s {
		v = append(v, n)
	}
	return StringArray{NewArray(v...)}
}

func (s *StringArray) Value() []string {
	return s.Strings()
}
