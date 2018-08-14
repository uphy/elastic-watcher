package actions

import (
	"github.com/uphy/elastic-watcher/watcher/context"
	"github.com/uphy/elastic-watcher/watcher/input"
)

type WebhookAction struct {
	input.HTTPRequest
}

func (w *WebhookAction) Run(ctx context.ExecutionContext) error {
	err := w.Execute(ctx)
	return err
}
