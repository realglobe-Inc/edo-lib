package driver

import (
	"testing"
)

func TestMongoIdProviderLister(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoIdProviderLister(mongoAddr, testLabel, "id-provider-lister", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*idProviderLister).base.(*mongoKeyValueStore).base.DB(testLabel).DropDatabase()

	if _, err := reg.(*idProviderLister).base.Put("list", testIdps); err != nil {
		t.Fatal(err)
	}

	testIdProviderLister(t, reg)
}

func TestMongoIdProviderListerStamp(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoIdProviderLister(mongoAddr, testLabel, "id-provider-lister", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*idProviderLister).base.(*mongoKeyValueStore).base.DB(testLabel).DropDatabase()

	if _, err := reg.(*idProviderLister).base.Put("list", testIdps); err != nil {
		t.Fatal(err)
	}

	testIdProviderListerStamp(t, reg)
}
