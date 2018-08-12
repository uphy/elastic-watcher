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
	if s.Lang != nil && *s.Lang != "javascript" {
		return nil, errors.New("unsupported script language: " + *s.Lang)
	}
	params := map[string]interface{}{}
	if s.Params != nil {
		m, err := s.Params.Map(ctx)
		if err != nil {
			return nil, err
		}
		for key, value := range m {
			params[key] = value
		}
	}
	var script string
	if s.Inline != nil {
		script = *s.Inline
	} else if s.Source != nil {
		script = *s.Source
	} else if s.ID != nil {
		var sid = *s.ID
		ss, ok := ctx.GlobalConfig().Scripts[sid]
		if !ok {
			return nil, errors.New("script not found: " + sid)
		}
		if ss.Lang != nil && *ss.Lang != "javascript" {
			return nil, errors.New("unsupported script language: " + *ss.Lang)
		}
		script = ss.Source
	}
	v, err := RunScript(ctx, script, params)
	if err != nil {
		return nil, err
	}
	return &v, err
}

func (i Script) MarshalJSON() ([]byte, error) {
	// full script
	if i.ID != nil || i.Lang != nil || i.Params != nil {
		type T Script
		var t T
		t = T(i)
		return json.Marshal(t)
	}
	// inline script
	if i.Inline != nil {
		return json.Marshal(i.Inline)
	} else if i.Source != nil {
		return json.Marshal(i.Source)
	}
	return nil, errors.New("empty script")
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
