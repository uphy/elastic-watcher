package condition

import "github.com/uphy/elastic-watcher/watcher/context"

type NeverCondition struct {
}

func (c NeverCondition) Match(ctx context.ExecutionContext) (bool, error) {
	return false, nil
}
