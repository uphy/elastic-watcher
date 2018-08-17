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
	ctx          context.ExecutionContext
}

func NewWatch(globalConfig *config.Config, c *WatchConfig) *Watch {
	return &Watch{c, globalConfig, context.New(globalConfig, c.Metadata)}
}

func (w *Watch) Run() error {
	ctx := w.ctx
	context.Init(ctx)

	runner := ctx.TaskRunner()

	// input
	if w.c.Input != nil {
		if err := runner.Run(w.c.Input); err != nil {
			return w.wrapError(ctx, "input", err)
		}
	}

	// check condition
	if w.c.Condition != nil {
		if err := runner.Run(condition.NewTask(w.c.Condition)); err != nil {
			return w.wrapError(ctx, "condition", err)
		}
	}

	// transform
	if w.c.Transform != nil {
		if err := runner.Run(w.c.Transform); err != nil {
			return w.wrapError(ctx, "transform", err)
		}
	}

	// run actions
	if w.c.Actions != nil {
		return w.wrapError(ctx, "action", runner.Run(w.c.Actions))
	}

	return nil
}

func (w *Watch) wrapError(ctx context.ExecutionContext, phase string, err error) error {
	if err == nil {
		return nil
	}
	return errors.Wrapf(err, "failed to run watch(id='%s') at '%s'", ctx.WatchID(), phase)
}
