package transform

import (
	"fmt"
	"testing"

	"github.com/uphy/elastic-watcher/watcher/context"
)

func TestScriptTransformer(t *testing.T) {
	ctx := context.TODO()
	script := `
	ctx.payload.foo = ['a','b'];
	ctx.payload
	`
	task := &ScriptTransform{
		context.Script{
			Inline: &script,
		},
	}
	if err := ctx.TaskRunner().Run(task); err != nil {
		t.Error(err)
	}
	fmt.Println(ctx.Payload())
}
