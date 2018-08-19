package transform

import "github.com/uphy/elastic-watcher/pkg/context"

type (
	SearchTransform struct {
		Request Request `json:"request"`
	}

	Request struct {
		Body    interface{} `json:"body"`
		Indices []string    `json:"indices"`
	}
)

func (s *SearchTransform) Run(ctx context.ExecutionContext) error {
	v, err := context.Search(ctx, s.Request.Indices, s.Request.Body)
	if err != nil {
		return err
	}
	ctx.SetPayload(v)
	return nil
}
