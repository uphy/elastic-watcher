package condition

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/uphy/elastic-watcher/watcher/context"
)

type (
	Condition interface {
		Match(ctx context.ExecutionContext) (bool, error)
	}
	Conditions struct {
		Script     *ScriptCondition  `json:"script,omitempty"`
		Compare    *CompareCondition `json:"compare,omitempty"`
		Always     *AlwaysCondition  `json:"always,omitempty"`
		Never      *NeverCondition   `json:"never,omitempty"`
		conditions []Condition       `json:"-"`
	}
	ConditionsTask struct {
		c Condition
	}
)

func NewTask(c Condition) context.Task {
	return &ConditionsTask{c}
}

func (c *ConditionsTask) Run(ctx context.ExecutionContext) error {
	matched, err := c.c.Match(ctx)
	if err != nil {
		return err
	}
	if !matched {
		return context.ErrStop
	}
	return nil
}

func (c *Conditions) Match(ctx context.ExecutionContext) (bool, error) {
	if c.conditions == nil {
		conditions := []Condition{}
		for _, condition := range []Condition{c.Compare, c.Script, c.Always, c.Never} {
			if reflect.ValueOf(condition).IsNil() {
				continue
			}
			conditions = append(conditions, condition)
		}
		c.conditions = conditions
	}

	for _, condition := range c.conditions {
		matched, err := condition.Match(ctx)
		if err != nil {
			return false, errors.Wrapf(err, "failed on condition")
		}
		if !matched {
			return false, nil
		}
	}
	return true, nil
}
