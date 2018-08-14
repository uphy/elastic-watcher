package transform

import (
	"github.com/uphy/elastic-watcher/watcher/context"
)

type ChainTransformer []Transform

func (c ChainTransformer) Run(ctx context.ExecutionContext) error {
	for _, t := range c {
		err := t.Run(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
