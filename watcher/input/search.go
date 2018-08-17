package input

import (
	"github.com/uphy/elastic-watcher/watcher/context"
	"github.com/uphy/elastic-watcher/watcher/transform"
)

type (
	SearchInput struct {
		transform.SearchTransform
	}
)

func (s SearchInput) Run(ctx context.ExecutionContext) error {
	return s.SearchTransform.Run(ctx)
}
