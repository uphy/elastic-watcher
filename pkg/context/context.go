package context

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/uphy/elastic-watcher/pkg/config"
)

type (
	ExecutionContext interface {
		ID() string
		WatchID() string
		ExecutionTime() time.Time
		Trigger() Trigger
		Metadata() JSONObject
		Vars() JSONObject
		SetVars(vars JSONObject)
		Payload() JSONObject
		SetPayload(payload JSONObject)
		GlobalConfig() *config.Config
		Logger() *logrus.Entry
		TaskRunner() *TaskRunner
	}
	JSONObject map[string]interface{}
	Trigger    struct {
		TriggeredTime time.Time `json:"triggered_time"`
		ScheduledTime time.Time `json:"scheduled_time"`
	}
)

func (j JSONObject) Clone() (JSONObject, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	dec := json.NewDecoder(&buf)
	err := enc.Encode(j)
	if err != nil {
		return nil, err
	}
	var copy map[string]interface{}
	err = dec.Decode(&copy)
	if err != nil {
		return nil, err
	}
	return copy, nil
}

func (j JSONObject) String() string {
	b, err := json.Marshal(j)
	if err != nil {
		return "<unabled to marshal JSONObject>"
	}
	return string(b)
}

func Init(ctx ExecutionContext) {
	ctx.SetPayload(JSONObject{})
	ctx.SetVars(JSONObject{})
	ctx.TaskRunner().Init()
}
