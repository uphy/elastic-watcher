package transform

import (
	"github.com/uphy/elastic-watcher/watcher/context"
)

type (
	Transform struct {
		Chain  *ChainTransformer  `json:"chain,omitempty"`
		Script *ScriptTransformer `json:"script,omitempty"`
		Search *SearchTransformer `json:"search,omitempty"`
	}
	Transformer interface {
		Transform(ctx context.ExecutionContext) error
	}
)

func (t *Transform) Transform(ctx context.ExecutionContext) error {
	var err error
	if t.Search != nil {
		err = t.Search.Transform(ctx)
		if err != nil {
			return err
		}
	}
	if t.Script != nil {
		err = t.Script.Transform(ctx)
		if err != nil {
			return err
		}
	}
	if t.Chain != nil {
		err = t.Chain.Transform(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
