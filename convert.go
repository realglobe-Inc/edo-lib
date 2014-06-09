package util

import (
	"github.com/idada/v8.go"
	"github.com/realglobe-Inc/edo/edoerror"
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

// JSON からデコードした interface{} を v8.Value にする。
func ToJsValue(engine *v8.Engine, param interface{}) (*v8.Value, error) {
	// case は "JSON and Go" の "Generic JSON with interface{}" を参考にした。
	switch p := param.(type) {
	case bool:
		return engine.NewBoolean(p), nil
	case float64:
		return engine.NewNumber(p), nil
	case string:
		return engine.NewString(p), nil
	case []interface{}:
		elems, err := ToJsArray(engine, p)
		if err != nil {
			return nil, erro.Wrap(err)
		}

		val := engine.NewArray(len(elems))
		array := val.ToArray()
		for i, elem := range elems {
			array.SetElement(i, elem)
		}
		return val, nil
	case map[string]interface{}:
		elems, err := ToJsMap(engine, p)
		if err != nil {
			return nil, erro.Wrap(err)
		}

		val := engine.NewObject()
		obj := val.ToObject()
		for key, elem := range elems {
			obj.SetProperty(key, elem, v8.PA_None)
		}
		return val, nil
	case nil:
		return engine.Null(), nil
	default:
		return nil, erro.Wrap(edoerror.NewInvalidParameter(param))
	}
}

// JSON からデコードした []interface{} を []*v8.Value にする。
func ToJsArray(engine *v8.Engine, param []interface{}) ([]*v8.Value, error) {
	vals := []*v8.Value{}
	for _, p := range param {
		val, err := ToJsValue(engine, p)
		if err != nil {
			return nil, erro.Wrap(err)
		}

		vals = append(vals, val)
	}
	return vals, nil
}

// JSON からデコードした map[string]interface{} を map[string]*v8.Value にする。
func ToJsMap(engine *v8.Engine, param map[string]interface{}) (map[string]*v8.Value, error) {
	vals := map[string]*v8.Value{}
	for key, p := range param {
		val, err := ToJsValue(engine, p)
		if err != nil {
			return nil, erro.Wrap(err)
		}

		vals[key] = val
	}
	return vals, nil
}
