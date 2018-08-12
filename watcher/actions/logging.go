package actions

import (
	"log"

	"github.com/uphy/elastic-watcher/watcher/context"
)

type LoggingAction struct {
	Text context.TemplateValue `json:"text"`
}

func (l *LoggingAction) Run(ctx context.ExecutionContext) error {
	s, err := l.Text.String(ctx)
	if err != nil {
		return err
	}
	log.Println(s)
	return nil
}
