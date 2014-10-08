package driver

import (
	"testing"
)

func TestMongoIdProviderAttributeRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoIdProviderAttributeRegistry(mongoAddr, testLabel, "id-provider-registry", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*idProviderAttributeRegistry).base.(*mongoKeyValueStore).base.DB(testLabel).DropDatabase()

	if _, err := reg.(*idProviderAttributeRegistry).base.Put(testIdpUuid+"/"+testAttrName, testAttr); err != nil {
		t.Fatal(err)
	}

	testIdProviderAttributeRegistry(t, reg)
}

func TestMongoIdProviderAttributeRegistryStamp(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoIdProviderAttributeRegistry(mongoAddr, testLabel, "id-provider-registry", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*idProviderAttributeRegistry).base.(*mongoKeyValueStore).base.DB(testLabel).DropDatabase()

	if _, err := reg.(*idProviderAttributeRegistry).base.Put(testIdpUuid+"/"+testAttrName, testAttr); err != nil {
		t.Fatal(err)
	}

	testIdProviderAttributeRegistryStamp(t, reg)
}
