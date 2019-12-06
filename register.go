package alejson

import "github.com/kode4food/ale/namespace"

import "github.com/kode4food/ale/data"

import "github.com/kode4food/ale/compiler/arity"

// Register alejson funcitons to the provided namespace
func Register(ns namespace.Type) {
	ns.Declare("marshal").Bind(data.MakeApplicative(AleMarshal, arity.MakeFixedChecker(1)))
	ns.Declare("unmarshal").Bind(data.MakeApplicative(AleUnmarshal, arity.MakeFixedChecker(1)))
}

// RegisterQualified register alejson package under the `json` package name
func RegisterQualified(m *namespace.Manager) namespace.Type {
	ns := m.GetQualified(data.Name("json"))
	Register(ns)
	return ns
}
