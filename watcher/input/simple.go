package input

import "github.com/uphy/elastic-watcher/watcher/context"

type SimpleInput struct {
	context.JSONObject
}

func (s SimpleInput) Run(ctx context.ExecutionContext) error {
	ctx.SetPayload(s.JSONObject)
	return nil
}
