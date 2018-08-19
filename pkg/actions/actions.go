package actions

import (
	"errors"
	"fmt"
	"reflect"

	"time"

	"github.com/uphy/elastic-watcher/pkg/condition"
	"github.com/uphy/elastic-watcher/pkg/context"
	"github.com/uphy/elastic-watcher/pkg/jsonutil"
	"github.com/uphy/elastic-watcher/pkg/transform"
)

type (
	Action interface {
		context.Task
	}
	DryRunner interface {
		DryRun(ctx context.ExecutionContext) error
	}
	actionContainer struct {
		// common options
		Condition      *condition.Conditions `json:"condition"`
		Transform      *transform.Transforms `json:"transform"`
		ThrottlePeriod *jsonutil.Duration    `json:"throttle_period"`
		DryRun         *bool                 `json:"dry_run"`

		// actions.  need to be specified one of these in json file.
		Logging   *LoggingAction   `json:"logging"`
		SendEmail *SendEmailAction `json:"send_email"`
		Webhook   *WebhookAction   `json:"webhook"`

		// lastAlert maps _key parameter of payload to last action time
		lastAlert map[string]time.Time `json:"-"`
		action    Action               `json:"-"`
	}

	Actions map[string]*actionContainer
)

func (a Actions) Run(ctx context.ExecutionContext) error {
	tasks := []context.Task{}
	for _, action := range a {
		tasks = append(tasks, context.NewTask(action.run))
	}
	return ctx.TaskRunner().Run(context.NewParallelTask(tasks))
}

func (a *actionContainer) run(ctx context.ExecutionContext) error {
	// initialization
	if a.action == nil {
		for _, action := range []Action{a.Logging, a.SendEmail, a.Webhook} {
			if reflect.ValueOf(action).IsNil() {
				continue
			}
			if a.action != nil {
				return errors.New("multiple actions specified in a name")
			}
			a.action = action
		}
		if a.action == nil {
			return errors.New("no action is defined")
		}
	}
	if a.lastAlert == nil {
		a.lastAlert = map[string]time.Time{}
	}

	// throttle
	keyObj := ctx.Payload()["_key"]
	key := fmt.Sprint(keyObj)
	if a.ThrottlePeriod != nil {
		if lastAlert, ok := a.lastAlert[key]; ok {
			if time.Now().Before(lastAlert.Add(a.ThrottlePeriod.Duration)) {
				return nil
			}
		}
	}

	// condition
	if a.Condition != nil {
		matched, err := a.Condition.Match(ctx)
		if err != nil {
			return err
		}
		if !matched {
			return nil
		}
	}

	// transform
	if a.Transform != nil {
		err := a.Transform.Run(ctx)
		if err != nil {
			return err
		}
	}

	// run action
	if a.DryRun != nil && *a.DryRun {
		if dryRunner, ok := a.action.(DryRunner); ok {
			if err := dryRunner.DryRun(ctx); err != nil {
				return err
			}
		} else {
			ctx.Logger().Infof("%s skipped because it doesn't support dry-run.", reflect.TypeOf(a.action))
		}
	} else {
		if err := a.action.Run(ctx); err != nil {
			return err
		}
	}

	a.lastAlert[key] = time.Now()
	return nil
}
