package context

import (
	"github.com/Sirupsen/logrus"
)

type scopedExecutionContext struct {
	ExecutionContext
	id         string
	logger     *logrus.Entry
	payload    JSONObject
	vars       JSONObject
	taskRunner *TaskRunner
}

func wrapContext(parent ExecutionContext, inheritTaskRunner bool) (ExecutionContext, error) {
	id := generateID()
	payload, err := parent.Payload().Clone()
	if err != nil {
		return nil, err
	}
	vars, err := parent.Vars().Clone()
	if err != nil {
		return nil, err
	}

	entry := parent.Logger()
	if parent.GlobalConfig().Debug {
		// overwrite id
		entry = entry.WithField("id", id)
	}
	wrapped := &scopedExecutionContext{
		ExecutionContext: parent,
		id:               id,
		payload:          payload,
		vars:             vars,
		logger:           entry,
	}
	wrapped.taskRunner = parent.TaskRunner()
	if !inheritTaskRunner {
		wrapped.taskRunner = parent.TaskRunner().addWorker(wrapped)
	}
	return wrapped, nil
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
func (s *scopedExecutionContext) TaskRunner() *TaskRunner {
	return s.taskRunner
}
