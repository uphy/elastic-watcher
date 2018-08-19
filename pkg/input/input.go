package input

import (
	"errors"
	"reflect"

	"github.com/uphy/elastic-watcher/pkg/condition"
	"github.com/uphy/elastic-watcher/pkg/context"
)

type Inputs struct {
	// common options
	Condition condition.Condition `json:"condition"`
	// input
	Chain     *ChainInput     `json:"chain"`
	Search    *SearchInput    `json:"search"`
	HTTP      *HTTPInput      `json:"http"`
	Simple    *SimpleInput    `json:"simple"`
	Transform *TransformInput `json:"transform"`
	input     Input           `json:"-"`
}

type Input interface {
	context.Task
}

func (i *Inputs) Run(ctx context.ExecutionContext) error {
	// initialization
	if i.input == nil {
		for _, input := range []Input{i.Chain, i.Search, i.HTTP, i.Simple, i.Transform} {
			if reflect.ValueOf(input).IsNil() {
				continue
			}
			if i.input != nil {
				return errors.New("multiple inputs specified in a name")
			}
			i.input = input
		}
		if i.input == nil {
			return errors.New("no input is defined")
		}
	}

	// condition
	if i.Condition != nil {
		matched, err := i.Condition.Match(ctx)
		if err != nil {
			return err
		}
		if !matched {
			ctx.Logger().Debug("Input has been skipped because condition is not matched.")
			return nil
		}
	}

	// input
	return i.input.Run(ctx)
}
