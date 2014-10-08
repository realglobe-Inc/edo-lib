package driver

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileIdProviderAttributeRegistry(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	reg := NewFileIdProviderAttributeRegistry(path, 0)
	if _, err := reg.(*idProviderAttributeRegistry).base.Put(testIdpUuid+"/"+testAttrName, testAttr); err != nil {
		t.Fatal(err)
	}

	testIdProviderAttributeRegistry(t, reg)
}

func TestFileIdProviderAttributeRegistryStamp(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	reg := NewFileIdProviderAttributeRegistry(path, 0)
	if _, err := reg.(*idProviderAttributeRegistry).base.Put(testIdpUuid+"/"+testAttrName, testAttr); err != nil {
		t.Fatal(err)
	}

	testIdProviderAttributeRegistryStamp(t, reg)
}
