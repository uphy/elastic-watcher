package condition

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/uphy/elastic-watcher/watcher/context"
)

type (
	Condition interface {
		Match(ctx context.ExecutionContext) (bool, error)
	}
	Conditions map[string]Condition
)

func (c Conditions) Match(ctx context.ExecutionContext) (bool, error) {
	for name, condition := range c {
		matched, err := condition.Match(ctx)
		if err != nil {
			return false, errors.Wrapf(err, "failed on condition '%s'", name)
		}
		if !matched {
			return false, nil
		}
	}
	return true, nil
}

func (c *Conditions) UnmarshalJSON(data []byte) (err error) {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	conditions := map[string]Condition{}
	for name, cond := range m {
		var condition Condition
		switch name {
		case "script":
			condition = &ScriptCondition{}
		case "compare":
			condition = &CompareCondition{}
		default:
			return errors.New("unsupported condition: " + name)
		}
		o, _ := json.Marshal(cond)
		if err := json.Unmarshal(o, condition); err != nil {
			return err
		}
		conditions[name] = condition
	}
	*c = conditions
	return nil
}
