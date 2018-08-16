package actions

import (
	"encoding/json"

	"github.com/uphy/elastic-watcher/watcher/context"
	"github.com/uphy/elastic-watcher/watcher/input"
)

type WebhookAction struct {
	input.HTTPRequest
}

func (w *WebhookAction) Run(ctx context.ExecutionContext) error {
	return w.Execute(ctx)
}

func (w *WebhookAction) DryRun(ctx context.ExecutionContext) error {
	b, err := json.Marshal(w.HTTPRequest)
	if err != nil {
		return err
	}
	ctx.Logger().Info("WebHook: " + string(b))
	return nil
}
