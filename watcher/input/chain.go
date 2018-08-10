package input

import (
	"github.com/uphy/elastic-watcher/watcher/context"
)

type Chain struct {
	Inputs []Input `json:"inputs"`
}

func (c *Chain) Read(ctx context.ExecutionContext) (interface{}, error) {
	for _, input := range c.Inputs {
		r, err := input.Read(ctx)
		if err != nil {
			return nil, err
		}
		ctx.SetPayload(r)
	}
	return ctx.Payload(), nil
}
