package util

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestStringSet(t *testing.T) {
	type testType struct {
		S map[string]StringSet
	}

	var a testType
	a.S = map[string]StringSet{"": NewStringSet(map[string]bool{"a": false, "b": true})}
	a.S[""]["c"] = true

	buff, err := json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}

	var b testType
	if err := json.Unmarshal(buff, &b); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(b, a) {
		t.Error(b)
	}
}

func TestMakeStringSet(t *testing.T) {
	m := map[string]bool{
		"a": true,
		"b": true,
		"c": true,
	}

	s1 := NewStringSet(m)
	if !reflect.DeepEqual(map[string]bool(s1), m) {
		t.Error(s1, m)
	}

	l := []string{}
	for elem := range m {
		l = append(l, elem)
	}
	s2 := StringSetFromSlice(l)

	if !reflect.DeepEqual(map[string]bool(s2), m) {
		t.Error(s2, m)
	}
}
