package transform

import (
	"fmt"

	"github.com/uphy/elastic-watcher/watcher/context"
)

type ScriptTransformer struct {
	context.Script
}

func (t *ScriptTransformer) Run(ctx context.ExecutionContext) error {
	v, err := t.Script.Value(ctx)
	if err != nil {
		return err
	}
	switch vv := v.(type) {
	case context.JSONObject:
		ctx.SetPayload(vv)
	case map[string]interface{}:
		ctx.SetPayload(vv)
	case nil:
		ctx.SetPayload(nil)
	default:
		return fmt.Errorf("incompatible return value. Script return value must be an javascript object. :%v", v)
	}
	return nil
}
