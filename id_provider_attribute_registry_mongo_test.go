package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMongoIdProviderAttributeRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoIdProviderAttributeRegistry(mongoAddr, testLabel, "id-provider-registry")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*idProviderRegistry).keyValueStore.(*mongoKeyValueStore).DB(testLabel).DropDatabase()

	if err := reg.(*idProviderRegistry).put(idProviderAttributeKey(testIdpUuid, testAttrName), testAttr); err != nil {
		t.Fatal(err)
	}

	testIdProviderAttributeRegistry(t, reg)
}

// キャッシュ用。
func TestMongoDatedIdProviderAttributeRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoDatedIdProviderAttributeRegistry(mongoAddr, testLabel, "id-provider-registry", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*datedIdProviderAttributeRegistry).datedKeyValueStore.(*mongoDatedKeyValueStore).DB(testLabel).DropDatabase()

	if _, err := reg.(*datedIdProviderAttributeRegistry).stampedPut(idProviderAttributeKey(testIdpUuid, testAttrName), testAttr); err != nil {
		t.Fatal(err)
	}

	testDatedIdProviderAttributeRegistry(t, reg)
}
