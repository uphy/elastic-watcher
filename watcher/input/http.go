package input

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/uphy/elastic-watcher/watcher/context"
)

type (
	HTTP struct {
		Request HTTPRequest `json:"request"`
	}
	HTTPRequest struct {
		// split url components
		Scheme *string                 `json:"scheme"`
		Host   *string                 `json:"host"`
		Port   *int                    `json:"port"`
		Path   *context.TemplateValue  `json:"path"`
		Params *context.TemplateValues `json:"params"`
		// simple url
		URL     *string                 `json:"url"`
		Method  *string                 `json:"method"`
		Body    *context.TemplateValue  `json:"body"`
		Headers *context.TemplateValues `json:"headers"`
	}
)

func (h *HTTP) Read(ctx context.ExecutionContext) (interface{}, error) {
	// url
	urlstring, err := h.buildURL(ctx)
	if err != nil {
		return nil, err
	}

	// request method
	method := "GET"
	if h.Request.Method != nil {
		method = *h.Request.Method
	}
	method = strings.ToUpper(method)

	// request body
	var body io.Reader
	if method != "GET" && h.Request.Body != nil {
		s, err := h.Request.Body.String(ctx)
		if err != nil {
			return nil, err
		}
		body = strings.NewReader(s)
	}

	// create http request
	req, err := http.NewRequest(method, urlstring, body)
	if err != nil {
		return nil, err
	}

	// add request headers
	if h.Request.Headers != nil {
		for _, key := range h.Request.Headers.Keys() {
			value, err := h.Request.Headers.StringByKey(ctx, key)
			if err != nil {
				return nil, err
			}
			req.Header.Add(key, value)
		}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	contentType := resp.Header.Get("Content-Type")

	payload := map[string]interface{}{}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	switch contentType {
	case "application/json", "text/json", "text/yaml", "text/x-yaml", "application/yaml", "application/x-yaml":
		if err := yaml.Unmarshal(respBody, &payload); err != nil {
			return nil, err
		}
	default:
		payload["_value"] = respBody
	}
	payloadHeaders := map[string]interface{}{}
	for k, v := range resp.Header {
		if len(v) == 1 {
			payloadHeaders[k] = v[0]
		} else {
			payloadHeaders[k] = v
		}
	}
	payload["_headers"] = payloadHeaders
	payload["_status_code"] = resp.StatusCode
	return payload, nil
}

func (h *HTTP) buildURL(ctx context.ExecutionContext) (string, error) {
	if h.Request.URL != nil {
		return *h.Request.URL, nil
	}
	if h.Request.Host == nil {
		return "", errors.New("host is required")
	}
	scheme := "http"
	if h.Request.Scheme != nil {
		scheme = *h.Request.Scheme
	}
	port := 80
	if h.Request.Port != nil {
		port = *h.Request.Port
	}
	path := ""
	if h.Request.Path != nil {
		p, err := h.Request.Path.String(ctx)
		if err != nil {
			return "", err
		}
		path = p
		if strings.HasPrefix(path, "/") {
			path = path[1:]
		}
	}
	queryParams := url.Values{}
	if h.Request.Params != nil {
		for _, key := range h.Request.Params.Keys() {
			v, err := h.Request.Params.StringByKey(ctx, key)
			if err != nil {
				return "", err
			}
			queryParams.Add(key, v)
		}
	}
	var q string
	if len(queryParams) > 0 {
		q = "?" + queryParams.Encode()
	}
	return fmt.Sprintf("%s://%s:%d/%s%s", scheme, *h.Request.Host, port, path, q), nil
}
