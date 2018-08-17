package input

import (
	"github.com/uphy/elastic-watcher/watcher/context"
)

type ChainInput struct {
	Inputs []NamedInput `json:"inputs"`
}

type NamedInput map[string]Inputs

func (i *ChainInput) Run(ctx context.ExecutionContext) error {
	p := context.JSONObject{}
	ctx.SetPayload(p)
	for _, input := range i.Inputs {
		for name, t := range input {
			if err := ctx.TaskRunner().Run(&t); err != nil {
				return err
			}
			payload, err := ctx.Payload().Clone()
			if err != nil {
				return err
			}
			p[name] = payload
			ctx.SetPayload(p)
		}
	}
	return nil
}
