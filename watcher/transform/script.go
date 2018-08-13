package transform

import (
	"github.com/uphy/elastic-watcher/watcher/context"
)

type ScriptTransformer struct {
	context.Script
}

func (t *ScriptTransformer) Transform(ctx context.ExecutionContext) error {
	v, err := t.Script.Value(ctx)
	if err != nil {
		return err
	}
	ctx.SetPayload(v)
	return nil
}
