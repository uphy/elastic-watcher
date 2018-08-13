package input

import (
	"github.com/uphy/elastic-watcher/watcher/context"
	"github.com/uphy/elastic-watcher/watcher/transform"
)

type TransformInput struct {
	transform.Transform
}

func (t *TransformInput) Read(ctx context.ExecutionContext) (context.Payload, error) {
	if err := t.Transform.Transform(ctx); err != nil {
		return nil, err
	}
	return ctx.Payload(), nil
}
