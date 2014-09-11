package driver

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileTimeLimitedKeyValueStore(t *testing.T) {
	path, err := ioutil.TempDir("", "test_edo_driver")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	testTimeLimitedKeyValueStore(t, NewFileTimeLimitedKeyValueStore(path))
}
