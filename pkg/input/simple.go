package input

import (
	"github.com/uphy/elastic-watcher/pkg/context"
)

type SimpleInput context.JSONObject

func (s SimpleInput) Run(ctx context.ExecutionContext) error {
	ctx.SetPayload(context.JSONObject(s))
	return nil
}
