package input

import "github.com/uphy/elastic-watcher/watcher/context"

type Simple map[string]interface{}

func (s *Simple) Read(ctx context.ExecutionContext) (interface{}, error) {
	return s, nil
}
