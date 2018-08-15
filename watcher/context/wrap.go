package context

import (
	"github.com/Sirupsen/logrus"
)

type scopedExecutionContext struct {
	ExecutionContext
	id      string
	logger  *logrus.Entry
	payload JSONObject
	vars    JSONObject
}

func newScopedContext(ctx ExecutionContext) (ExecutionContext, error) {
	id := generateID()
	payload, err := ctx.Payload().Clone()
	if err != nil {
		return nil, err
	}
	vars, err := ctx.Vars().Clone()
	if err != nil {
		return nil, err
	}

	entry := ctx.Logger()
	if ctx.GlobalConfig().Debug {
		// overwrite id
		entry = entry.WithField("id", id)
	}
	return &scopedExecutionContext{
		ExecutionContext: ctx,
		id:               id,
		payload:          payload,
		vars:             vars,
		logger:           entry,
	}, nil
}
func (s *scopedExecutionContext) ID() string {
	return s.id
}
func (s *scopedExecutionContext) Vars() JSONObject {
	return s.vars
}
func (s *scopedExecutionContext) SetVars(vars JSONObject) {
	s.vars = vars
}
func (s *scopedExecutionContext) Payload() JSONObject {
	return s.payload
}
func (s *scopedExecutionContext) SetPayload(payload JSONObject) {
	s.payload = payload
}
func (s *scopedExecutionContext) Logger() *logrus.Entry {
	return s.logger
}
