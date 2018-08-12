package context

import (
	"encoding/json"
	"errors"

	"github.com/robertkrimen/otto"
)

type Script struct {
	Inline *string        `json:"inline,omitempty"`
	Source *string        `json:"source,omitempty"`
	ID     *string        `json:"id,omitempty"`
	Lang   *string        `json:"lang,omitempty"`
	Params TemplateValues `json:"params,omitempty"`
}

func (s *Script) Bool(ctx ExecutionContext) (bool, error) {
	v, err := s.value(ctx)
	if err != nil {
		return false, err
	}
	return v.ToBoolean()
}

func (s *Script) Value(ctx ExecutionContext) (interface{}, error) {
	v, err := s.value(ctx)
	if err != nil {
		return false, err
	}
	return v.Export()
}

func (s *Script) value(ctx ExecutionContext) (*otto.Value, error) {
	if s.Lang != nil && *s.Lang == "painless" {
		return nil, errors.New("`painless` script is not supported")
	}
	var params map[string]interface{}
	if s.Params != nil {
		m, err := s.Params.Map(ctx)
		if err != nil {
			return nil, err
		}
		for key, value := range m {
			params[key] = value
		}
	}
	script, err := s.findScript()
	if err != nil {
		return nil, err
	}
	v, err := RunScript(ctx, script, params)
	if err != nil {
		return nil, err
	}
	return &v, err
}

func (s *Script) findScript() (string, error) {
	if s.Inline != nil {
		return *s.Inline, nil
	}
	if s.Source != nil {
		return *s.Source, nil
	}
	if s.ID != nil {
		return "", nil
	}
	return "", errors.New("empty script")
}

func (i Script) MarshalJSON() ([]byte, error) {
	// full script
	if i.ID != nil || i.Lang != nil || i.Params != nil {
		type T Script
		var t T
		return json.Marshal(t)
	}
	// inline script
	script, err := i.findScript()
	if err != nil {
		return nil, err
	}
	return json.Marshal(script)
}

func (i *Script) UnmarshalJSON(data []byte) error {
	// full script
	type T Script
	var t T
	if err := json.Unmarshal(data, &t); err == nil {
		*i = Script(t)
		return nil
	}
	// inline script
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*i = Script{
		Source: &s,
	}
	return nil
}
