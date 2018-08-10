package context

import (
	"time"

	"github.com/uphy/elastic-watcher/config"
)

type scopedExecutionContext struct {
	ctx     ExecutionContext
	payload interface{}
	vars    interface{}
}

func newScopedContext(ctx ExecutionContext) ExecutionContext {
	return &scopedExecutionContext{
		ctx:     ctx,
		payload: ctx.Payload(),
		vars:    ctx.Vars(),
	}
}

func (s *scopedExecutionContext) WatchID() string {
	return s.ctx.WatchID()
}
func (s *scopedExecutionContext) ExecutionTime() time.Time {
	return s.ctx.ExecutionTime()
}
func (s *scopedExecutionContext) Trigger() Trigger {
	return s.ctx.Trigger()
}
func (s *scopedExecutionContext) Metadata() interface{} {
	return s.ctx.Metadata()
}
func (s *scopedExecutionContext) Vars() interface{} {
	return s.vars
}
func (s *scopedExecutionContext) SetVars(vars interface{}) {
	s.vars = vars
}
func (s *scopedExecutionContext) Payload() interface{} {
	return s.payload
}
func (s *scopedExecutionContext) SetPayload(payload interface{}) {
	s.payload = payload
}
func (s *scopedExecutionContext) GlobalConfig() *config.Config {
	return s.ctx.GlobalConfig()
}
