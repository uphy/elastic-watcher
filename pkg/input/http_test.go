package input

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/uphy/elastic-watcher/pkg/context"
)

func TestHTTPRun(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/json", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Add("Content-Type", "application/json")
		resp.WriteHeader(200)
		resp.Write([]byte(`{"foo":"bar"}`))
	})

	server := &http.Server{Addr: ":12345", Handler: mux}
	go func() {
		server.ListenAndServe()
	}()
	defer server.Close()
	time.Sleep(time.Second)

	req := `
	{
		"request": {
			"scheme": "http",
			"host": "localhost",
			"port": 12345,
			"method": "GET",
			"path": "/json"
		}
	}
	`
	var httpInput HTTPInput
	if err := json.Unmarshal([]byte(req), &httpInput); err != nil {
		t.Error(err)
	}
	ctx := context.TODO()
	err := httpInput.Run(ctx)
	if err != nil {
		t.Error(err)
	}
	vv := ctx.Payload()
	if vv["_status_code"] != 200 {
		t.Errorf("want 200 but %v", vv["_status_code"])
	}
	if vv["foo"] != "bar" {
		t.Errorf("want bar but %v", vv["foo"])
	}
}
