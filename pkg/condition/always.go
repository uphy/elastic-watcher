package condition

import "github.com/uphy/elastic-watcher/pkg/context"

type AlwaysCondition struct {
}

func (c AlwaysCondition) Match(ctx context.ExecutionContext) (bool, error) {
	return true, nil
}
