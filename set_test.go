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
