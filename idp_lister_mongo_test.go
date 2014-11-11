package driver

import (
	"testing"
)

func TestMongoIdpLister(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoIdpLister(mongoAddr, testLabel, "id-provider-lister", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*idpLister).base.(*mongoKeyValueStore).base.DB(testLabel).DropDatabase()

	if _, err := reg.(*idpLister).base.Put("list", testIdps); err != nil {
		t.Fatal(err)
	}

	testIdpLister(t, reg)
}

func TestMongoIdpListerStamp(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoIdpLister(mongoAddr, testLabel, "id-provider-lister", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*idpLister).base.(*mongoKeyValueStore).base.DB(testLabel).DropDatabase()

	if _, err := reg.(*idpLister).base.Put("list", testIdps); err != nil {
		t.Fatal(err)
	}

	testIdpListerStamp(t, reg)
}
