package input

import (
	"github.com/uphy/elastic-watcher/watcher/context"
	"github.com/uphy/elastic-watcher/watcher/transform"
)

type TransformInput struct {
	transform.Transforms
}

func (t *TransformInput) Run(ctx context.ExecutionContext) error {
	return t.Transforms.Run(ctx)
}
