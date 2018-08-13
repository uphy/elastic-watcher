package input

import (
	"github.com/uphy/elastic-watcher/watcher/context"
	"github.com/uphy/elastic-watcher/watcher/transform"
)

type (
	SearchInput struct {
		transform.SearchTransformer
	}
)

func (s SearchInput) Read(ctx context.ExecutionContext) (context.Payload, error) {
	if err := s.Transform(ctx); err != nil {
		return nil, err
	}
	return ctx.Payload(), nil
}
