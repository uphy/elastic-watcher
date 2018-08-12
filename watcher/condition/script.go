package condition

import (
	"github.com/uphy/elastic-watcher/watcher/context"
)

type ScriptCondition struct {
	context.Script
}

func (s *ScriptCondition) Match(ctx context.ExecutionContext) (bool, error) {
	return s.Bool(ctx)
}
