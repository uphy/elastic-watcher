package input

import (
	"encoding/json"
	"errors"

	"github.com/uphy/elastic-watcher/watcher/context"
)

type Input struct {
	typ    string
	reader Reader
}

type Reader interface {
	context.Task
}

func (i *Input) Run(ctx context.ExecutionContext) error {
	return i.reader.Run(ctx)
}

func (i Input) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(i.reader)
	if err != nil {
		return nil, err
	}
	var reader interface{}
	if err := json.Unmarshal(j, &reader); err != nil {
		return nil, err
	}
	v := map[string]interface{}{}
	v[i.typ] = reader
	return json.Marshal(v)
}

func (i *Input) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if len(v) == 0 {
		return errors.New("empty input")
	} else if len(v) > 1 {
		return errors.New("multiple input not supported")
	}

	for name, r := range v {
		var reader Reader
		switch name {
		case "search":
			reader = &SearchInput{}
		case "simple":
			reader = &SimpleInput{}
		case "transform":
			reader = &TransformInput{}
		case "http":
			reader = &HTTPInput{}
		case "chain":
			reader = &ChainInput{}
		default:
			return errors.New("unsupported input: " + name)
		}

		j, err := json.Marshal(r)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(j, reader); err != nil {
			return err
		}
		i.typ = name
		i.reader = reader
	}
	return nil
}
