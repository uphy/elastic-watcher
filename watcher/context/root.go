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
		metadata      JSONObject
		payload       JSONObject
		vars          JSONObject
		globalConfig  *config.Config
		logger        *logrus.Entry
		taskRunner    *TaskRunner
	}
)

func TODO() ExecutionContext {
	return New(&config.Config{}, JSONObject{})
}

func New(globalConfig *config.Config, metadata JSONObject) ExecutionContext {
	id := generateID()
	watchID := generateID()
	t := time.Now()

	// logger
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}
	if globalConfig.Debug {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}
	// logger entry
	entry := logrus.NewEntry(logger)
	if globalConfig.Debug {
		entry = entry.WithFields(logrus.Fields{
			"watchID": watchID,
			"id":      id,
		})
	}

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
		payload:      JSONObject{},
		vars:         JSONObject{},
		logger:       entry,
	}

	runner := NewTaskRunner(ctx)
	ctx.taskRunner = runner
	return ctx
}

func Init(ctx ExecutionContext) {
	ctx.SetPayload(JSONObject{})
	ctx.SetVars(JSONObject{})
	ctx.TaskRunner().Init()
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

func (e *rootExecutionContext) Metadata() JSONObject {
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

func (r *rootExecutionContext) Logger() *logrus.Entry {
	return r.logger
}

func (r *rootExecutionContext) TaskRunner() *TaskRunner {
	return r.taskRunner
}
