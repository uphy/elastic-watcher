package input

import (
	"encoding/json"
	"errors"

	"github.com/uphy/elastic-watcher/watcher/condition"
	"github.com/uphy/elastic-watcher/watcher/context"
)

type Input struct {
	typ       string
	reader    Reader
	condition condition.Condition
}

type Reader interface {
	context.Task
}

func (i *Input) Run(ctx context.ExecutionContext) error {
	if i.condition != nil {
		matched, err := i.condition.Match(ctx)
		if err != nil {
			return err
		}
		if !matched {
			ctx.Logger().Debug("Input has been skipped because condition is not matched.")
			return nil
		}
	}
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
	}

	for name, r := range v {
		var reader Reader
		var cond condition.Condition
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
		case "condition":
			cond = &condition.Conditions{}
		default:
			return errors.New("unsupported input: " + name)
		}

		j, err := json.Marshal(r)
		if err != nil {
			return err
		}

		if reader != nil {
			if err := json.Unmarshal(j, reader); err != nil {
				return err
			}
			i.typ = name
			if i.reader != nil {
				return errors.New("multiple input not supported")
			}
			i.reader = reader
		} else if cond != nil {
			if err := json.Unmarshal(j, cond); err != nil {
				return err
			}
			i.condition = cond
		}
	}
	return nil
}
