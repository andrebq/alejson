package alejson

import (
	"github.com/kode4food/ale/data"
	"github.com/valyala/fastjson"
	"math"
)

// AleUnmarshal wraps Unmarshal in a Ale compatible function
func AleUnmarshal(args ...data.Value) data.Value {
	v, ok := args[0].(data.String)
	if !ok {
		v = data.String(args[0].String())
	}
	ret, err := Unmarshal(string(v))
	if err != nil {
		return data.Null
	}
	return ret
}

// Unmarshal input into an ale Value following this rules:
//
// JSON Array -> data.Vector
// JSON Object -> data.Object
// JSON Object Keys -> data.Keyword
// Int values within int32 range -> data.Integer
// Int values outside of int32 range -> data.Float
// Float values -> data.Float
// True/False -> data.True/data.False
// String -> String
func Unmarshal(input string) (data.Value, error) {
	value, err := fastjson.Parse(input)
	if err != nil {
		return data.Null, err
	}
	return parseValue(value), nil
}

func parseValue(p *fastjson.Value) data.Value {
	switch t := p.Type(); t {
	case fastjson.TypeArray:
		return parseArray(p)
	case fastjson.TypeObject:
		return parseObject(p)
	case fastjson.TypeString:
		b, _ := p.StringBytes()
		return data.String(b)
	case fastjson.TypeFalse:
		return data.False
	case fastjson.TypeTrue:
		return data.True
	case fastjson.TypeNumber:
		fval, _ := p.Float64()
		if fval > math.MaxInt32 || fval < math.MinInt32 {
			return data.Float(fval)
		}
		if fval == float64(int32(fval)) {
			return data.Integer(int(fval))
		}
		return data.Float(fval)
	}
	panic("Unrecognized type: " + p.Type().String())
}

func parseArray(v *fastjson.Value) data.Sequence {
	inputVals, _ := v.Array()
	vals := make([]data.Value, len(inputVals))
	for i, v := range inputVals {
		vals[i] = parseValue(v)
	}
	return data.NewVector(vals...)
}

func parseObject(v *fastjson.Value) data.Object {
	inputObj, _ := v.Object()
	outputObj := data.NewObject()
	inputObj.Visit(func(k []byte, v *fastjson.Value) {
		outputObj[data.Keyword(string(k))] = parseValue(v)
	})
	return outputObj
}
