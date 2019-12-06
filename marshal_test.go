package alejson

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/core/bootstrap"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/eval"
)

func TestMarshalAllTypes(t *testing.T) {
	manager := bootstrap.TopLevelManager()
	bootstrap.Into(manager)
	RegisterQualified(manager)

	type testCase struct {
		aleCode string
		output  string
	}
	tests := []testCase{
		{`{:a "hello"}`, `{"a":"hello"}`},
		{`{:a 1234}`, `{"a":1234}`},
		{`{:a 1234.5}`, `{"a":1234.5}`},
		{`{:a true}`, `{"a":true}`},
		{`{:a false}`, `{"a":false}`},
		{`{:a {:b {:c "hi"}}}`, `{"a":{"b":{"c":"hi"}}}`},
		{`{:a [1 2 false true :abc "cde" {:a 1}]}`, `{"a":[1,2,false,true,"abc","cde",{"a":1}]}`},
		{`{:a (list 1 2 false true :abc "cde" {:a 1})}`, `{"a":[1,2,false,true,"abc","cde",{"a":1}]}`},
	}
	for _, tc := range tests {
		ns := manager.GetAnonymous()
		output := eval.String(ns, data.String(fmt.Sprintf(`(json/marshal %v)`, tc.aleCode)))
		if str, ok := output.(data.String); !ok {
			t.Errorf("Result of marshal-json is not data.String: %v", output)
		} else if string(str) != tc.output {
			t.Errorf("Invalid json string %v", str)
		}
	}
}
