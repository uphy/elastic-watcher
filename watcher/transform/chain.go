package transform

import (
	"github.com/uphy/elastic-watcher/watcher/context"
)

type ChainTransformer []Transform

func (c ChainTransformer) Transform(ctx context.ExecutionContext) error {
	for _, t := range c {
		err := t.Transform(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
