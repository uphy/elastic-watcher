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
		context.Task
	}
)

func (t *Transform) Run(ctx context.ExecutionContext) error {
	transformers := []Transformer{}
	if t.Chain != nil {
		transformers = append(transformers, t.Chain)
	}
	if t.Script != nil {
		transformers = append(transformers, t.Script)
	}
	if t.Search != nil {
		transformers = append(transformers, t.Search)
	}
	for _, transformer := range transformers {
		if err := transformer.Run(ctx); err != nil {
			return err
		}
	}
	return nil
}
