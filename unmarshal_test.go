package alejson

import (
	"bytes"
	"testing"

	"github.com/kode4food/ale/core/bootstrap"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/eval"
	"github.com/kode4food/ale/stdlib"
)

func TestUnmarshalAllTypes(t *testing.T) {
	manager := bootstrap.TopLevelManager()
	bootstrap.Into(manager)
	RegisterQualified(manager)

	type testCase struct {
		ale  string
		json string
	}
	tests := []testCase{
		{`{:a "hello"}`, `{"a":"hello"}`},
		{`{:a 1234}`, `{"a":1234}`},
		{`{:a 1234.5}`, `{"a":1234.5}`},
		{`{:a true}`, `{"a":true}`},
		{`{:a false}`, `{"a":false}`},
		{`{:a {:b {:c "hi"}}}`, `{"a":{"b":{"c":"hi"}}}`},
		{`{:a [1 2 false true "abc" "cde" {:a 1}]}`, `{"a":[1,2,false,true,"abc","cde",{"a":1}]}`},
	}
	for _, tc := range tests {
		ns := manager.GetAnonymous()
		ns.Declare(data.Name("input")).Bind(data.String(tc.json))
		output := eval.String(ns, data.String(`(json/unmarshal input)`))
		if str := stdOutput(output); str != tc.ale {
			t.Errorf("Invalid ale string %v should be %s", str, tc.ale)
		}
	}
}

func stdOutput(v data.Value) string {
	buf := &bytes.Buffer{}
	writer := stdlib.NewWriter(buf, stdlib.StrOutput)
	writer.Write(v)
	return buf.String()
}
