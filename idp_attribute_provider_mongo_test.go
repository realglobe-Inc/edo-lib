package driver

import (
	"testing"
)

func TestMongoIdpAttributeProvider(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoIdpAttributeProvider(mongoAddr, testLabel, "idp-attribute-provider", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*idpAttributeProvider).base.(*mongoKeyValueStore).base.DB(testLabel).DropDatabase()

	if _, err := reg.(*idpAttributeProvider).base.Put(testIdpUuid+"/"+testAttrName, testAttr); err != nil {
		t.Fatal(err)
	}

	testIdpAttributeProvider(t, reg)
}

func TestMongoIdpAttributeProviderStamp(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoIdpAttributeProvider(mongoAddr, testLabel, "idp-attribute-provider", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*idpAttributeProvider).base.(*mongoKeyValueStore).base.DB(testLabel).DropDatabase()

	if _, err := reg.(*idpAttributeProvider).base.Put(testIdpUuid+"/"+testAttrName, testAttr); err != nil {
		t.Fatal(err)
	}

	testIdpAttributeProviderStamp(t, reg)
}
