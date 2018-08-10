package input

import "github.com/uphy/elastic-watcher/watcher/context"

type (
	Search struct {
		Request Request `json:"request"`
	}

	Request struct {
		Body    interface{} `json:"body"`
		Indices []string    `json:"indices"`
	}
)

func (s Search) Read(ctx context.ExecutionContext) (interface{}, error) {
	return context.Search(ctx, s.Request.Indices, s.Request.Body)
}
