package actions

import (
	"encoding/json"

	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/uphy/elastic-watcher/watcher/condition"
	"github.com/uphy/elastic-watcher/watcher/context"
	"github.com/uphy/elastic-watcher/watcher/transform"
)

type (
	Action interface {
		Run(ctx context.ExecutionContext) error
	}

	actionContainer struct {
		typ            string
		action         Action
		condition      *condition.Conditions
		transform      *transform.Transform
		throttlePeriod *Duration
		lastAlert      time.Time
	}

	Actions map[string]*actionContainer
)

func (a *actionContainer) run(ctx context.ExecutionContext) error {
	if a.throttlePeriod != nil {
		if time.Now().Before(a.lastAlert.Add(a.throttlePeriod.Duration)) {
			return nil
		}
	}
	if a.condition != nil {
		matched, err := a.condition.Match(ctx)
		if err != nil {
			return err
		}
		if !matched {
			return nil
		}
	}
	if a.transform != nil {
		v, err := a.transform.Read(ctx)
		if err != nil {
			return err
		}
		ctx.SetPayload(v)
	}
	if err := a.action.Run(ctx); err != nil {
		return err
	}
	a.lastAlert = time.Now()
	return nil
}

func (a Actions) Run(ctx context.ExecutionContext) error {
	var errs error
	for name, action := range a {
		if err := action.run(context.Wrap(ctx)); err != nil {
			errs = multierror.Append(errs, errors.Wrap(err, "failed to execute action: "+name))
		}
	}
	return errs
}

func (a *Actions) UnmarshalJSON(data []byte) (err error) {
	var m map[string](map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	actions := map[string]*actionContainer{}
	for name, typeAndOptions := range m {
		a := &actionContainer{}
		for typ, options := range typeAndOptions {
			b, err := json.Marshal(options)
			if err != nil {
				return err
			}

			switch typ {
			case "condition":
				var cond condition.Conditions
				if err := json.Unmarshal(b, &cond); err != nil {
					return err
				}
				a.condition = &cond
			case "throttle_period":
				var d Duration
				if err := json.Unmarshal(b, &d); err != nil {
					return err
				}
				a.throttlePeriod = &d
			case "transform":
				var t transform.Transform
				if err := json.Unmarshal(b, &t); err != nil {
					return err
				}
				a.transform = &t
			default:
				act := newAction(typ)
				if act == nil {
					return errors.New("unsupported action: " + typ)
				}
				if err := json.Unmarshal(b, act); err != nil {
					return err
				}
				a.typ = typ
				a.action = act
			}
		}
		actions[name] = a
	}
	*a = actions
	return nil
}

func (a Actions) MarshalJSON() ([]byte, error) {
	m := map[string](map[string]interface{}){}
	for name, action := range a {
		actionMap := map[string]interface{}{}
		d, err := json.Marshal(action.action)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(d, &actionMap); err != nil {
			return nil, err
		}
		if action.condition != nil {
			b, err := json.Marshal(action.condition)
			if err != nil {
				return nil, err
			}
			conditionMap := map[string]interface{}{}
			if err := json.Unmarshal(b, &conditionMap); err != nil {
				return nil, err
			}
			actionMap["condition"] = conditionMap
		}
		if action.throttlePeriod != nil {
			actionMap["throttle_period"] = action.throttlePeriod.String()
		}
		if action.transform != nil {
			b, err := json.Marshal(action.transform)
			if err != nil {
				return nil, err
			}
			transformMap := map[string]interface{}{}
			if err := json.Unmarshal(b, &transformMap); err != nil {
				return nil, err
			}
			actionMap["transform"] = transformMap
		}
		m[name] = actionMap
	}
	return json.Marshal(m)
}

func newAction(typ string) Action {
	switch typ {
	case "logging":
		return &LoggingAction{}
	case "send_email":
		return &SendEmailAction{}
	default:
		return nil
	}
}
