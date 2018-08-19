package transform

import (
	"github.com/uphy/elastic-watcher/pkg/context"
)

type (
	Transforms struct {
		Chain  *ChainTransform  `json:"chain,omitempty"`
		Script *ScriptTransform `json:"script,omitempty"`
		Search *SearchTransform `json:"search,omitempty"`
	}
	Transform interface {
		context.Task
	}
)

func (t *Transforms) Run(ctx context.ExecutionContext) error {
	transformers := []Transform{}
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
