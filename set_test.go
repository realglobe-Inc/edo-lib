package util

import (
	"encoding/json"
	"testing"
)

func TestStringSet(t *testing.T) {
	type testType struct {
		S map[string]*StringSet
	}

	var a testType
	a.S = map[string]*StringSet{"": NewStringSet(map[string]bool{"a": false, "b": true})}
	a.S[""].Put("c")

	buff, err := json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}

	var b testType
	if err := json.Unmarshal(buff, &b); err != nil {
		t.Fatal(err)
	}

	if len(b.S) != 1 {
		t.Error(b.S)
	} else if b.S[""] == nil {
		t.Error(b.S)
	} else if b.S[""].Contains("a") {
		t.Error(b.S[""])
	} else if !b.S[""].Contains("b") {
		t.Error(b.S[""])
	} else if !b.S[""].Contains("c") {
		t.Error(b.S[""])
	}

	b.S[""].Remove("c")
	if b.S[""].Contains("c") {
		t.Error(b.S[""])
	}
}
