package util

import (
	"github.com/idada/v8.go"
	"testing"
)

func TestBoolToJsValue(t *testing.T) {
	v8.NewEngine().NewContext(nil).Scope(func(cs v8.ContextScope) {
		val, err := ToJsValue(cs.GetEngine(), true)
		if err != nil {
			t.Fatal(err)
		}
		if !val.IsBoolean() || val.ToBoolean() != true {
			t.Error(string(v8.ToJSON(val)))
		}
	})
}

func TestNumberToJsValue(t *testing.T) {
	v8.NewEngine().NewContext(nil).Scope(func(cs v8.ContextScope) {
		val, err := ToJsValue(cs.GetEngine(), 1.234)
		if err != nil {
			t.Fatal(err)
		}
		if !val.IsNumber() || val.ToNumber() != 1.234 {
			t.Error(string(v8.ToJSON(val)))
		}
	})
}

func TestStringToJsValue(t *testing.T) {
	v8.NewEngine().NewContext(nil).Scope(func(cs v8.ContextScope) {
		val, err := ToJsValue(cs.GetEngine(), "aaaa")
		if err != nil {
			t.Fatal(err)
		}
		if !val.IsString() || val.ToString() != "aaaa" {
			t.Error(string(v8.ToJSON(val)))
		}
	})
}

func TestSliceToJsValue(t *testing.T) {
	v8.NewEngine().NewContext(nil).Scope(func(cs v8.ContextScope) {
		val, err := ToJsValue(cs.GetEngine(), []interface{}{"a", "bb", "ccc"})
		if err != nil {
			t.Fatal(err)
		}
		if !val.IsArray() {
			t.Error(string(v8.ToJSON(val)))
		} else if ary := val.ToArray(); ary.Length() != 3 || ary.GetElement(0).ToString() != "a" || ary.GetElement(1).ToString() != "bb" || ary.GetElement(2).ToString() != "ccc" {
			t.Error(string(v8.ToJSON(val)))
		}
	})
}

func TestMapToJsValue(t *testing.T) {
	v8.NewEngine().NewContext(nil).Scope(func(cs v8.ContextScope) {
		val, err := ToJsValue(cs.GetEngine(), map[string]interface{}{"A": "a", "B": "bb", "C": "ccc"})
		if err != nil {
			t.Fatal(err)
		}
		if !val.IsObject() {
			t.Error(string(v8.ToJSON(val)))
		} else if obj := val.ToObject(); !obj.HasProperty("A") || obj.GetProperty("A").ToString() != "a" || !obj.HasProperty("B") || obj.GetProperty("B").ToString() != "bb" || !obj.HasProperty("C") || obj.GetProperty("C").ToString() != "ccc" {
			t.Error(string(v8.ToJSON(val)))
		}
	})
}
