package input

import (
	"errors"

	"github.com/uphy/elastic-watcher/watcher/context"
)

type Transform struct {
	Script *context.Script `json:"script"`
}

func (t *Transform) Read(ctx context.ExecutionContext) (interface{}, error) {
	if t.Script == nil {
		return nil, errors.New("`script` not defined at `transform`")
	}
	return t.Script.Value(ctx)
}
