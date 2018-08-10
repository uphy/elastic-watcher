package input

import "github.com/uphy/elastic-watcher/watcher/context"

type Transform struct {
	Script string `json:"script"`
}

func (t *Transform) Read(ctx context.ExecutionContext) (interface{}, error) {
	v, err := context.RunScript(ctx, t.Script)
	if err != nil {
		return nil, err
	}
	return v.Export()
}
