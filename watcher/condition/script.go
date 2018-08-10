package condition

import (
	"errors"

	"github.com/uphy/elastic-watcher/watcher/context"
)

type ScriptCondition struct {
	Source string                 `json:"source"`
	Lang   string                 `json:"lang"`
	Params map[string]interface{} `json:"params,omitempty"`
}

func (s *ScriptCondition) Match(ctx context.ExecutionContext) (bool, error) {
	if s.Lang != "javascript" {
		return false, errors.New("unsupported lang: " + s.Lang)
	}
	v, err := context.RunScript(ctx, s.Source)
	if err != nil {
		return false, err
	}
	return v.ToBoolean()
}
