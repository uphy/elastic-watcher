package context

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/cbroglie/mustache"
	"github.com/robertkrimen/otto"
)

func data(ctx ExecutionContext) (interface{}, error) {
	v := map[string]interface{}{}
	v["watch_id"] = ctx.WatchID()
	v["execution_time"] = ctx.ExecutionTime()
	v["trigger"] = map[string]interface{}{
		"triggered_time": ctx.Trigger().TriggeredTime,
		"scheduled_time": ctx.Trigger().ScheduledTime,
	}
	v["metadata"] = ctx.Metadata()
	v["payload"] = ctx.Payload()
	v["vars"] = ctx.Vars()
	return v, nil
}

func RenderTemplate(ctx ExecutionContext, template string) (string, error) {
	data, err := data(ctx)
	if err != nil {
		return "", err
	}
	return mustache.Render(template, map[string]interface{}{"ctx": data})
}

func RunScript(ctx ExecutionContext, script string) (otto.Value, error) {
	// initialize javascript engine
	vm := otto.New()

	// set context
	ctxData, err := data(ctx)
	if err != nil {
		return otto.NullValue(), err
	}
	if err := vm.Set("ctx", ctxData); err != nil {
		return otto.NullValue(), err
	}
	// run the script
	v, err := vm.Run(script)
	if err != nil {
		return otto.NullValue(), err
	}
	return v, nil
}

func Search(ctx ExecutionContext, indices []string, query interface{}) (interface{}, error) {
	client := &http.Client{
		Transport: http.DefaultTransport,
	}
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	encodedIndices := []string{}
	for _, index := range indices {
		encodedIndices = append(encodedIndices, url.QueryEscape(index))
	}
	resp, err := client.Post(fmt.Sprintf("%s/%s/_search", ctx.GlobalConfig().Elasticsearch.URL, strings.Join(encodedIndices, ",")), "application/json", bytes.NewReader(queryJSON))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to search. (status=%d, body=%s)", resp.StatusCode, string(b))
	}
	d := json.NewDecoder(resp.Body)
	var m map[string]interface{}
	if err := d.Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}