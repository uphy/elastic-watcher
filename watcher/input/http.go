package input

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/uphy/elastic-watcher/watcher/context"
)

type (
	HTTPInput struct {
		Request HTTPRequest `json:"request"`
	}
	HTTPRequest struct {
		// split url components
		Scheme *string                 `json:"scheme,omitempty"`
		Host   *string                 `json:"host,omitempty"`
		Port   *int                    `json:"port,omitempty"`
		Path   *context.TemplateValue  `json:"path,omitempty"`
		Params *context.TemplateValues `json:"params,omitempty"`
		// simple url
		URL     *string                 `json:"url,omitempty"`
		Method  *string                 `json:"method,omitempty"`
		Body    *context.TemplateValue  `json:"body,omitempty"`
		Headers *context.TemplateValues `json:"headers,omitempty"`
	}
)

func (h *HTTPInput) Read(ctx context.ExecutionContext) (context.Payload, error) {
	return h.Request.Execute(ctx)
}
func (h *HTTPRequest) Execute(ctx context.ExecutionContext) (context.Payload, error) {
	// url
	urlstring, err := h.buildURL(ctx)
	if err != nil {
		return nil, err
	}

	// request method
	method := "GET"
	if h.Method != nil {
		method = *h.Method
	}
	method = strings.ToUpper(method)

	// request body
	var body io.Reader
	if method != "GET" && h.Body != nil {
		s, err := h.Body.String(ctx)
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
	if h.Headers != nil {
		for _, key := range h.Headers.Keys() {
			value, err := h.Headers.StringByKey(ctx, key)
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
	contentType, _, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))

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

func (h *HTTPRequest) buildURL(ctx context.ExecutionContext) (string, error) {
	if h.URL != nil {
		return *h.URL, nil
	}
	if h.Host == nil {
		return "", errors.New("host is required")
	}
	scheme := "http"
	if h.Scheme != nil {
		scheme = *h.Scheme
	}
	port := 80
	if h.Port != nil {
		port = *h.Port
	}
	path := ""
	if h.Path != nil {
		p, err := h.Path.String(ctx)
		if err != nil {
			return "", err
		}
		path = p
		if strings.HasPrefix(path, "/") {
			path = path[1:]
		}
	}
	queryParams := url.Values{}
	if h.Params != nil {
		for _, key := range h.Params.Keys() {
			v, err := h.Params.StringByKey(ctx, key)
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
	return fmt.Sprintf("%s://%s:%d/%s%s", scheme, *h.Host, port, path, q), nil
}
