package context

import (
	"fmt"
	"time"

	"github.com/uphy/elastic-watcher/config"
)

type (
	rootExecutionContext struct {
		watchID       string      `json:"watch_id"`
		executionTime time.Time   `json:"execution_time"`
		trigger       Trigger     `json:"trigger"`
		metadata      interface{} `json:"metadata"`
		payload       interface{} `json:"-"`
		vars          interface{} `json:"-"`
		globalConfig  *config.Config
	}
)

var currentID = 0

func New(globalConfig *config.Config) ExecutionContext {
	id := currentID
	currentID++
	t := time.Now()
	return &rootExecutionContext{
		watchID:       fmt.Sprint(id),
		executionTime: t,
		trigger: Trigger{
			TriggeredTime: t,
			ScheduledTime: t,
		},
		globalConfig: globalConfig,
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

func (e *rootExecutionContext) Payload() interface{} {
	return e.payload
}

func (e *rootExecutionContext) SetPayload(payload interface{}) {
	e.payload = payload
}

func (e *rootExecutionContext) GlobalConfig() *config.Config {
	return e.globalConfig
}
