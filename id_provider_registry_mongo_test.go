package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMongoIdProviderRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoIdProviderRegistry(mongoAddr, "test_driver_mongo", "idp")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*idProviderRegistry).keyValueStore.(*mongoKeyValueStore).DB("test_driver_mongo").DropDatabase()

	if err := reg.(*idProviderRegistry).put("a_b-c", "https://localhost:1234/query"); err != nil {
		t.Fatal(err)
	}

	testIdProviderRegistry(t, reg)
}

// キャッシュ用。
func TestMongoDatedIdProviderRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoDatedIdProviderRegistry(mongoAddr, "test_driver_mongo", "idp", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*datedIdProviderRegistry).datedKeyValueStore.(*mongoDatedKeyValueStore).DB("test_driver_mongo").DropDatabase()

	if _, err := reg.(*datedIdProviderRegistry).stampedPut("a_b-c", "https://localhost:1234/query"); err != nil {
		t.Fatal(err)
	}

	testDatedIdProviderRegistry(t, reg)
}
