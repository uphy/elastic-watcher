package context

import (
	"encoding/json"
	"testing"
)

func TestTemplateUnmarshal_SingleValue(t *testing.T) {
	type Foo struct {
		V TemplateValues `json:"v"`
	}
	var foo Foo
	j := `{
		"v": "{{ ctx.payload.value }}"
	}`
	if err := json.Unmarshal([]byte(j), &foo); err != nil {
		t.Error(err)
	}
	if foo.V.Size() != 1 {
		t.Errorf("want 1 but %d", foo.V.Size())
	}
	ctx := TODO()
	ctx.SetPayload(JSONObject{
		"value": "foo",
	})
	v, err := foo.V.String(ctx, 0)
	if err != nil {
		t.Error(err)
	}
	if v != "foo" {
		t.Errorf("want foo but %s", v)
	}
	b, err := json.Marshal(foo)
	if err != nil {
		t.Error(err)
	}
	var foo2 Foo
	if err := json.Unmarshal(b, &foo2); err != nil {
		t.Error(err)
	}
	if len(foo.V) != len(foo2.V) {
		t.Fail()
	}
	for _, key := range foo.V.Keys() {
		v1, err := foo.V.StringByKey(ctx, key)
		if err != nil {
			t.Error(err)
		}
		v2, err := foo2.V.StringByKey(ctx, key)
		if err != nil {
			t.Error(err)
		}
		if v1 != v2 {
			t.Fail()
		}
	}
}

func TestTemplateUnmarshal_ArrayValue(t *testing.T) {
	type Foo struct {
		V TemplateValues `json:"v"`
	}
	var foo Foo
	j := `{
		"v": ["{{ ctx.payload.value1 }}","{{ ctx.payload.value2 }}"]
	}`
	if err := json.Unmarshal([]byte(j), &foo); err != nil {
		t.Error(err)
	}
	if foo.V.Size() != 2 {
		t.Errorf("want 2 but %d", foo.V.Size())
	}
	ctx := TODO()
	ctx.SetPayload(JSONObject{
		"value1": "foo",
		"value2": "bar",
	})
	v1, err := foo.V.String(ctx, 0)
	if err != nil {
		t.Error(err)
	}
	if v1 != "foo" {
		t.Errorf("want foo but %s", v1)
	}
	v2, err := foo.V.String(ctx, 1)
	if err != nil {
		t.Error(err)
	}
	if v2 != "bar" {
		t.Errorf("want bar but %s", v1)
	}
}

func TestTemplateUnmarshal_MapValue(t *testing.T) {
	type Foo struct {
		V TemplateValues `json:"v"`
	}
	var foo Foo
	j := `{
		"v": {
			"key1": "{{ ctx.payload.value1 }}",
			"key2": "{{ ctx.payload.value2 }}"
		}
	}`
	if err := json.Unmarshal([]byte(j), &foo); err != nil {
		t.Error(err)
	}
	if foo.V.Size() != 2 {
		t.Errorf("want 2 but %d", foo.V.Size())
	}
	ctx := TODO()
	ctx.SetPayload(JSONObject{
		"value1": "foo",
		"value2": "bar",
	})
	v1, err := foo.V.StringByKey(ctx, "key1")
	if err != nil {
		t.Error(err)
	}
	if v1 != "foo" {
		t.Errorf("want foo but %s", v1)
	}
	v2, err := foo.V.StringByKey(ctx, "key2")
	if err != nil {
		t.Error(err)
	}
	if v2 != "bar" {
		t.Errorf("want bar but %s", v1)
	}
}
