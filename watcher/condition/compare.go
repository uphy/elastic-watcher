package condition

import (
	"encoding/json"
	"fmt"

	"github.com/uphy/elastic-watcher/watcher/context"
)

type (
	CompareCondition map[string]Compare

	Compare struct {
		EQ    *json.Number `json:"eq,omitempty"`
		NotEQ *json.Number `json:"not_eq,omitempty"`
		GT    *json.Number `json:"gt,omitempty"`
		GTE   *json.Number `json:"gte,omitempty"`
		LT    *json.Number `json:"lt,omitempty"`
		LTE   *json.Number `json:"lte,omitempty"`
	}
)

func (c CompareCondition) Match(ctx context.ExecutionContext) (bool, error) {
	for field, compare := range c {
		matched, err := compare.Match(field, ctx)
		if err != nil {
			return false, err
		}
		if !matched {
			return false, nil
		}
	}
	return true, nil
}

func (c *Compare) match(ctx context.ExecutionContext, format string, field string, value json.Number) (bool, error) {
	v, err := context.RunScript(ctx, fmt.Sprintf(format, field, value))
	if err != nil {
		return false, err
	}
	b, err := v.ToBoolean()
	if err != nil {
		return false, err
	}
	if !b {
		return false, nil
	}
	return true, nil
}

func (c *Compare) Match(field string, ctx context.ExecutionContext) (bool, error) {
	if c.EQ != nil {
		v, err := c.match(ctx, "%s === %v", field, *c.EQ)
		if err != nil {
			return false, err
		}
		if !v {
			return false, nil
		}
	}
	if c.NotEQ != nil {
		v, err := c.match(ctx, "%s !== %v", field, *c.NotEQ)
		if err != nil {
			return false, err
		}
		if !v {
			return false, nil
		}
	}
	if c.GTE != nil {
		v, err := c.match(ctx, "%s >= %v", field, *c.GTE)
		if err != nil {
			return false, err
		}
		if !v {
			return false, nil
		}
	}
	if c.GT != nil {
		v, err := c.match(ctx, "%s > %v", field, *c.GT)
		if err != nil {
			return false, err
		}
		if !v {
			return false, nil
		}
	}
	if c.LTE != nil {
		v, err := c.match(ctx, "%s <= %v", field, *c.LTE)
		if err != nil {
			return false, err
		}
		if !v {
			return false, nil
		}
	}
	if c.LT != nil {
		v, err := c.match(ctx, "%s < %v", field, *c.LT)
		if err != nil {
			return false, err
		}
		if !v {
			return false, nil
		}
	}
	return true, nil
}
