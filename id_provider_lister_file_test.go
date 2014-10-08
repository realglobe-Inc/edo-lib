package driver

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileIdProviderLister(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	reg := NewFileIdProviderLister(path, 0)
	if _, err := reg.(*idProviderLister).base.Put("list", testIdps); err != nil {
		t.Fatal(err)
	}

	testIdProviderLister(t, reg)
}

func TestFileIdProviderListerStamp(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	reg := NewFileIdProviderLister(path, 0)
	if _, err := reg.(*idProviderLister).base.Put("list", testIdps); err != nil {
		t.Fatal(err)
	}

	testIdProviderListerStamp(t, reg)
}
