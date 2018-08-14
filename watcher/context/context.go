package context

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/uphy/elastic-watcher/config"
)

type (
	ExecutionContext interface {
		WatchID() string
		ExecutionTime() time.Time
		Trigger() Trigger
		Metadata() interface{}
		Vars() interface{}
		SetVars(vars interface{})
		Payload() Payload
		SetPayload(payload Payload)
		GlobalConfig() *config.Config
		Logger() *logrus.Logger
	}
	Payload interface{}
	Trigger struct {
		TriggeredTime time.Time `json:"triggered_time"`
		ScheduledTime time.Time `json:"scheduled_time"`
	}
)

func Wrap(ctx ExecutionContext) ExecutionContext {
	return newScopedContext(ctx)
}
