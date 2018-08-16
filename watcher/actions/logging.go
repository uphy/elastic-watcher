package actions

import (
	"strings"

	"github.com/uphy/elastic-watcher/watcher/context"
)

type LoggingAction struct {
	Text context.TemplateValue `json:"text"`
}

func (l *LoggingAction) DryRun(ctx context.ExecutionContext) error {
	return l.Run(ctx)
}

func (l *LoggingAction) Run(ctx context.ExecutionContext) error {
	s, err := l.Text.String(ctx)
	if err != nil {
		return err
	}
	for _, ss := range strings.Split(s, "\n") {
		ctx.Logger().Info(ss)
	}
	return nil
}
