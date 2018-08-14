package input

import (
	"github.com/uphy/elastic-watcher/watcher/context"
	"github.com/uphy/elastic-watcher/watcher/transform"
)

type TransformInput struct {
	transform.Transform
}

func (t *TransformInput) Run(ctx context.ExecutionContext) error {
	return t.Transform.Run(ctx)
}
