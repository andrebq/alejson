package alejson

import (
	"fmt"
	"math"

	"github.com/kode4food/ale/data"
	"github.com/valyala/fastjson"
)

// AleMarshal wraps Marshal in a Ale compatible function
func AleMarshal(args ...data.Value) data.Value {
	buf, err := Marshal(args[0])
	if err != nil {
		return data.Null
	}
	return data.String(string(buf))
}

// Marshal encodes a data.Value to a JSON literal
//
// data.Sequence other than Object -> JSON Array
// data.Object -> JSON Object
// data.Integer within int32 range -> JSON Integer
// data.Integer outside int32 range -> JSON Float (without decimal)
// data.Float -> JSON Float
// data.String -> JSON String
// data.Keywordd -> JSON String
// data.True/data.False -> JSON True/False
func Marshal(v data.Value) ([]byte, error) {
	a := &fastjson.Arena{}
	mval, err := marshalValue(a, v)
	if err != nil {
		return nil, err
	}
	a.Reset()
	return mval.MarshalTo(nil), nil
}

func marshalValue(a *fastjson.Arena, v data.Value) (*fastjson.Value, error) {
	switch v := v.(type) {
	case data.Object:
		return marshalObject(a, v)
	case data.String:
		return a.NewString(string(v)), nil
	case data.Keyword:
		return a.NewString(string(v)), nil
	case data.Integer:
		if v > math.MaxInt32 || v < math.MaxInt32 {
			return a.NewNumberFloat64(float64(v)), nil
		}
		return a.NewNumberInt(int(v)), nil
	case data.Float:
		return a.NewNumberFloat64(float64(v)), nil
	case data.Bool:
		if v {
			return a.NewTrue(), nil
		}
		return a.NewFalse(), nil
	case data.Sequence:
		return marshalArray(a, v)
	}
	return nil, fmt.Errorf("alejson: not supported %#v", v)
}

func marshalObject(a *fastjson.Arena, aleObj data.Object) (*fastjson.Value, error) {
	obj := a.NewObject()
	for k, v := range aleObj {
		kw, ok := k.(data.Keyword)
		if !ok {
			return nil, fmt.Errorf("alejson: object keys MUST be keywords")
		}
		vv, err := marshalValue(a, v)
		if err != nil {
			return nil, err
		}
		obj.Set(string(kw), vv)
	}
	return obj, nil
}

func marshalArray(a *fastjson.Arena, v data.Sequence) (*fastjson.Value, error) {
	ret := a.NewArray()
	var idx int
	item := v.First()
	v = v.Rest()
	for item != data.Null {
		mv, err := marshalValue(a, item)
		if err != nil {
			return nil, err
		}
		ret.SetArrayItem(idx, mv)
		idx++
		item = v.First()
		v = v.Rest()
	}
	return ret, nil
}
