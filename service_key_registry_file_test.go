package driver

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileServiceKeyRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	reg := NewFileServiceKeyRegistry(path, 0)
	if _, err := reg.(*serviceKeyRegistry).base.Put(testServUuid, testPublicKey); err != nil {
		t.Fatal(err)
	}

	testServiceKeyRegistry(t, reg)
}

func TestFileServiceKeyRegistryStamp(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	reg := NewFileServiceKeyRegistry(path, 0)
	if _, err := reg.(*serviceKeyRegistry).base.Put(testServUuid, testPublicKey); err != nil {
		t.Fatal(err)
	}

	testServiceKeyRegistry(t, reg)
}
