package driver

import (
	"testing"
)

func testLister(t *testing.T, drv Lister) {
	defer drv.Close()

	keys, _, err := drv.Keys(nil)
	if err != nil {
		t.Fatal(err)
	} else if len(keys) != 1 || !keys[testKey] {
		t.Error(keys)
	}
}
