package watcher

import (
	"github.com/pkg/errors"
	"github.com/uphy/elastic-watcher/config"
	"github.com/uphy/elastic-watcher/watcher/condition"
	"github.com/uphy/elastic-watcher/watcher/context"
)

type Watch struct {
	c            *WatchConfig
	globalConfig *config.Config
}

func NewWatch(globalConfig *config.Config, c *WatchConfig) *Watch {
	return &Watch{c, globalConfig}
}

func (w *Watch) Run() error {
	// clear state
	ctx := context.New(w.globalConfig, w.c.Metadata)
	runner := ctx.TaskRunner()

	// input
	if err := runner.Run(&w.c.Input); err != nil {
		return w.wrapError(ctx, "input", err)
	}

	// check condition
	if err := runner.Run(condition.NewTask(w.c.Condition)); err != nil {
		return w.wrapError(ctx, "condition", err)
	}

	// transform
	if w.c.Transform != nil {
		if err := runner.Run(w.c.Transform); err != nil {
			return w.wrapError(ctx, "transform", err)
		}
	}

	// run actions
	return w.wrapError(ctx, "action", runner.Run(w.c.Actions))
}

func (w *Watch) wrapError(ctx context.ExecutionContext, phase string, err error) error {
	if err == nil {
		return nil
	}
	return errors.Wrapf(err, "failed to run watch(id='%s') at '%s'", ctx.WatchID(), phase)
}
