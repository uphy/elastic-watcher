package transform

import (
	"github.com/uphy/elastic-watcher/watcher/context"
)

type ChainTransform []Transforms

func (c ChainTransform) Run(ctx context.ExecutionContext) error {
	for _, t := range c {
		err := t.Run(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
