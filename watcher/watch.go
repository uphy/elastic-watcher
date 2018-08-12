package watcher

import (
	"github.com/pkg/errors"
	"github.com/uphy/elastic-watcher/config"
	"github.com/uphy/elastic-watcher/watcher/context"
)

type Watch struct {
	c   *WatchConfig
	ctx context.ExecutionContext
}

func NewWatch(globalConfig *config.Config, c *WatchConfig) *Watch {
	return &Watch{c, context.New(globalConfig, c.Metadata)}
}

func (w *Watch) Run() error {
	// clear state
	w.ctx.SetPayload(nil)

	// input
	data, err := w.c.Input.Read(w.ctx)
	if err != nil {
		return w.wrapError("input", err)
	}
	w.ctx.SetPayload(data)

	// check condition
	matched, err := w.c.Condition.Match(w.ctx)
	if err != nil {
		return w.wrapError("condition", err)
	}
	if !matched {
		return nil
	}

	// transform
	if w.c.Transform != nil {
		data, err := w.c.Transform.Read(w.ctx)
		if err != nil {
			return w.wrapError("transform", err)
		}
		w.ctx.SetPayload(data)
	}

	// run actions
	return w.wrapError("action", w.c.Actions.Run(w.ctx))
}

func (w *Watch) wrapError(phase string, err error) error {
	if err == nil {
		return nil
	}
	return errors.Wrapf(err, "failed to run watch(id='%s') at '%s'", w.ctx.WatchID(), phase)
}
