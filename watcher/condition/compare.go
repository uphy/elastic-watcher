package condition

import (
	"fmt"

	"github.com/uphy/elastic-watcher/watcher/context"
)

type (
	CompareCondition map[string]Compare

	Compare struct {
		EQ    *context.TemplateValue `json:"eq,omitempty"`
		NotEQ *context.TemplateValue `json:"not_eq,omitempty"`
		GT    *context.TemplateValue `json:"gt,omitempty"`
		GTE   *context.TemplateValue `json:"gte,omitempty"`
		LT    *context.TemplateValue `json:"lt,omitempty"`
		LTE   *context.TemplateValue `json:"lte,omitempty"`
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

func (c *Compare) match(ctx context.ExecutionContext, format string, field string, value string) (bool, error) {
	v, err := context.RunScript(ctx, fmt.Sprintf(format, field, value), nil)
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
		s, err := c.EQ.String(ctx)
		if err != nil {
			return false, err
		}
		v, err := c.match(ctx, "%s === '%v'", field, s)
		if err != nil {
			return false, err
		}
		if !v {
			return false, nil
		}
	}
	if c.NotEQ != nil {
		s, err := c.NotEQ.String(ctx)
		if err != nil {
			return false, err
		}
		v, err := c.match(ctx, "%s !== '%v'", field, s)
		if err != nil {
			return false, err
		}
		if !v {
			return false, nil
		}
	}
	if c.GTE != nil {
		s, err := c.GT.String(ctx)
		if err != nil {
			return false, err
		}
		v, err := c.match(ctx, "%s >= %v", field, s)
		if err != nil {
			return false, err
		}
		if !v {
			return false, nil
		}
	}
	if c.GT != nil {
		s, err := c.GT.String(ctx)
		if err != nil {
			return false, err
		}
		v, err := c.match(ctx, "%s > %v", field, s)
		if err != nil {
			return false, err
		}
		if !v {
			return false, nil
		}
	}
	if c.LTE != nil {
		s, err := c.LTE.String(ctx)
		if err != nil {
			return false, err
		}
		v, err := c.match(ctx, "%s <= %v", field, s)
		if err != nil {
			return false, err
		}
		if !v {
			return false, nil
		}
	}
	if c.LT != nil {
		s, err := c.LT.String(ctx)
		if err != nil {
			return false, err
		}
		v, err := c.match(ctx, "%s < %v", field, s)
		if err != nil {
			return false, err
		}
		if !v {
			return false, nil
		}
	}
	return true, nil
}
