package input

import (
	"github.com/uphy/elastic-watcher/watcher/context"
)

type (
	HTTP struct {
		Request HTTPRequest `json:"request"`
	}
	HTTPRequest struct {
		Scheme  *string
		Host    string
		Port    int
		Path    *string
		Method  *string
		Headers map[string]string
	}
)

func (h *HTTP) Read(ctx context.ExecutionContext) (interface{}, error) {
	return nil, nil
}
