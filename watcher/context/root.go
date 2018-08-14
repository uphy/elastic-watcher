package context

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/uphy/elastic-watcher/config"
)

type (
	rootExecutionContext struct {
		id            string
		watchID       string
		executionTime time.Time
		trigger       Trigger
		metadata      map[string]interface{}
		payload       JSONObject
		vars          JSONObject
		globalConfig  *config.Config
		logger        *logrus.Logger
		taskRunner    *TaskRunner
	}
)

func TODO() ExecutionContext {
	return New(&config.Config{}, map[string]interface{}{})
}

func New(globalConfig *config.Config, metadata map[string]interface{}) ExecutionContext {
	t := time.Now()
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	ctx := &rootExecutionContext{
		id:            generateID(),
		watchID:       generateID(),
		executionTime: t,
		trigger: Trigger{
			TriggeredTime: t,
			ScheduledTime: t,
		},
		metadata:     metadata,
		globalConfig: globalConfig,
		logger:       logger,
		payload:      JSONObject{},
		vars:         JSONObject{},
	}
	runner := NewTaskRunner(ctx)
	ctx.taskRunner = runner
	return ctx
}

func (e *rootExecutionContext) ID() string {
	return e.id
}
func (e *rootExecutionContext) WatchID() string {
	return e.watchID
}

func (e *rootExecutionContext) ExecutionTime() time.Time {
	return e.executionTime
}

func (e *rootExecutionContext) Trigger() Trigger {
	return e.trigger
}

func (e *rootExecutionContext) Metadata() interface{} {
	return e.metadata
}

func (e *rootExecutionContext) Vars() JSONObject {
	return e.vars
}

func (e *rootExecutionContext) SetVars(vars JSONObject) {
	e.vars = vars
}

func (e *rootExecutionContext) Payload() JSONObject {
	return e.payload
}

func (e *rootExecutionContext) SetPayload(payload JSONObject) {
	e.payload = payload
}

func (e *rootExecutionContext) GlobalConfig() *config.Config {
	return e.globalConfig
}

func (r *rootExecutionContext) Logger() *logrus.Logger {
	return r.logger
}

func (r *rootExecutionContext) TaskRunner() *TaskRunner {
	return r.taskRunner
}
