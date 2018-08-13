package input

import "github.com/uphy/elastic-watcher/watcher/context"

type SimpleInput struct {
	context.Payload
}

func (s SimpleInput) Read(ctx context.ExecutionContext) (context.Payload, error) {
	return s.Payload, nil
}
