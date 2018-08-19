package context

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/frohmut/mustache"
	"github.com/robertkrimen/otto"
	"github.com/rs/xid"
)

func data(ctx ExecutionContext) (map[string]interface{}, error) {
	v := map[string]interface{}{}
	v["watch_id"] = ctx.WatchID()
	v["execution_time"] = ctx.ExecutionTime()
	v["trigger"] = map[string]interface{}{
		"triggered_time": ctx.Trigger().TriggeredTime,
		"scheduled_time": ctx.Trigger().ScheduledTime,
	}
	v["metadata"] = ctx.Metadata()
	v["payload"] = map[string]interface{}(ctx.Payload())
	v["vars"] = ctx.Vars()
	return v, nil
}

func renderTemplate(ctx ExecutionContext, template string, params map[string]interface{}) (string, error) {
	data, err := data(ctx)
	if err != nil {
		return "", err
	}
	p := JSONObject{}
	p["ctx"] = data
	if params != nil {
		for k, v := range params {
			p[k] = v
		}
	}

	// https://github.com/elastic/elasticsearch/blob/master/docs/reference/search/search-template.asciidoc
	p["toJson"] = func(text string, render mustache.RenderFn) (string, error) {
		value, err := RunScript(ctx, text, nil)
		if err != nil {
			return "", err
		}
		v, err := value.Export()
		if err != nil {
			return "", err
		}
		b, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	return mustache.Render(template, p)
}

func setSplittedPayload(ctx ExecutionContext, v []JSONObject) {
	ctx.Vars()["__splitted__"] = v
}

func consumeSplittedPayload(ctx ExecutionContext) []JSONObject {
	vars := ctx.Vars()
	splitted, exist := vars["__splitted__"]
	if !exist {
		return nil
	}
	casted := splitted.([]JSONObject)
	delete(vars, "__splitted__")
	return casted
}

func RunScript(ctx ExecutionContext, script string, params map[string]interface{}) (otto.Value, error) {
	if ctx.GlobalConfig().Debug {
		ctx.Logger().Debug("Running script")
		ctx.Logger().Debugf("Script:\n%s", script)
		ctx.Logger().Debugf("Params:%s", params)
	}
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
	if err := vm.Set("log", func(call otto.FunctionCall) otto.Value {
		ctx.Logger().Info(call.Argument(0))
		return otto.UndefinedValue()
	}); err != nil {
		return otto.NullValue(), err
	}
	if err := vm.Set("split", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 {
			panic(call.Otto.MakeRangeError("splitted payload is required."))
		}
		input := call.Argument(0)
		s, _ := input.Export()

		items := reflect.ValueOf(s)
		if items.Kind() != reflect.Slice {
			panic(call.Otto.MakeTypeError(fmt.Sprintf("splitted payload must be an object array: %v", reflect.TypeOf(s))))
		}
		jsonObjects := []JSONObject{}
		for i := 0; i < items.Len(); i++ {
			elm := items.Index(i).Interface()
			casted, ok := elm.(map[string]interface{})
			if !ok {
				panic(call.Otto.MakeTypeError(fmt.Sprintf("element of splitted payload must be an object: %v", reflect.TypeOf(s))))
			}
			jsonObjects = append(jsonObjects, JSONObject(casted))
		}
		setSplittedPayload(ctx, jsonObjects)
		return otto.UndefinedValue()
	}); err != nil {
		return otto.NullValue(), err
	}
	if params != nil {
		for k, v := range params {
			if err := vm.Set(k, v); err != nil {
				return otto.NullValue(), err
			}
		}
	}

	// run the script
	v, err := vm.Run(script)
	if err != nil {
		return otto.NullValue(), err
	}
	return v, nil
}

func Search(ctx ExecutionContext, indices []string, query interface{}) (map[string]interface{}, error) {
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

func generateID() string {
	return xid.New().String()
}
