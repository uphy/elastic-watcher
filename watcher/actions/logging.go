package actions

import (
	"log"

	"github.com/uphy/elastic-watcher/watcher/context"
)

type LoggingAction struct {
	Text string `json:"text"`
}

func (l *LoggingAction) Run(ctx context.ExecutionContext) error {
	rendered, err := context.RenderTemplate(ctx, l.Text)
	if err != nil {
		return err
	}
	log.Println(rendered)
	return nil
}
