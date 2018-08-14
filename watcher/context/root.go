package context

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/uphy/elastic-watcher/config"
)

type (
	rootExecutionContext struct {
		watchID       string
		executionTime time.Time
		trigger       Trigger
		metadata      map[string]interface{}
		payload       Payload
		vars          interface{}
		globalConfig  *config.Config
		logger        *logrus.Logger
	}
)

var currentID = 0

func TODO() ExecutionContext {
	return New(&config.Config{}, map[string]interface{}{})
}

func New(globalConfig *config.Config, metadata map[string]interface{}) ExecutionContext {
	id := currentID
	currentID++
	t := time.Now()
	logger := logrus.New()

	return &rootExecutionContext{
		watchID:       fmt.Sprint(id),
		executionTime: t,
		trigger: Trigger{
			TriggeredTime: t,
			ScheduledTime: t,
		},
		metadata:     metadata,
		globalConfig: globalConfig,
		logger:       logger,
	}
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

func (e *rootExecutionContext) Vars() interface{} {
	return e.vars
}

func (e *rootExecutionContext) SetVars(vars interface{}) {
	e.vars = vars
}

func (e *rootExecutionContext) Payload() Payload {
	return e.payload
}

func (e *rootExecutionContext) SetPayload(payload Payload) {
	e.payload = payload
}

func (e *rootExecutionContext) GlobalConfig() *config.Config {
	return e.globalConfig
}

func (r *rootExecutionContext) Logger() *logrus.Logger {
	return r.logger
}
