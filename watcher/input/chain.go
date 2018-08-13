package input

import (
	"github.com/uphy/elastic-watcher/watcher/context"
)

type ChainInput struct {
	Inputs []NamedInput `json:"inputs"`
}

type NamedInput map[string]Input

func (i *ChainInput) Read(ctx context.ExecutionContext) (context.Payload, error) {
	p := map[string]interface{}{}
	for _, input := range i.Inputs {
		for name, t := range input {
			ctxForInput := context.Wrap(ctx)
			v, err := t.Read(ctxForInput)
			if err != nil {
				return nil, err
			}
			p[name] = v
			ctx.SetPayload(p)
		}
	}
	return ctx.Payload(), nil
}
