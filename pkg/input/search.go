package input

import (
	"github.com/uphy/elastic-watcher/pkg/context"
	"github.com/uphy/elastic-watcher/pkg/transform"
)

type (
	SearchInput struct {
		transform.SearchTransform
	}
)

func (s SearchInput) Run(ctx context.ExecutionContext) error {
	return s.SearchTransform.Run(ctx)
}
