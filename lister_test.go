package driver

import (
	"testing"
)

func testLister(t *testing.T, reg Lister) {
	keys, _, err := reg.Keys(nil)
	if err != nil {
		t.Fatal(err)
	} else if len(keys) != 1 || !keys[testKey] {
		t.Error(keys)
	}
}
