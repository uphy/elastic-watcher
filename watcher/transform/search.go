package transform

import "github.com/uphy/elastic-watcher/watcher/context"

type (
	SearchTransformer struct {
		Request Request `json:"request"`
	}

	Request struct {
		Body    interface{} `json:"body"`
		Indices []string    `json:"indices"`
	}
)

func (s SearchTransformer) Transform(ctx context.ExecutionContext) error {
	v, err := context.Search(ctx, s.Request.Indices, s.Request.Body)
	if err != nil {
		return err
	}
	ctx.SetPayload(v)
	return nil
}
