package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMongoServiceKeyRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoServiceKeyRegistry(mongoAddr, "test_driver_mongo", "key")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*serviceKeyRegistry).keyValueStore.(*mongoKeyValueStore).DB("test_driver_mongo").DropDatabase()

	if err := reg.(*serviceKeyRegistry).put("a_b-c", "kore ga kagi dayo."); err != nil {
		t.Fatal(err)
	}

	testServiceKeyRegistry(t, reg)
}

// キャッシュ用。
func TestMongoDatedServiceKeyRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoDatedServiceKeyRegistry(mongoAddr, "test_driver_mongo", "key", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*datedServiceKeyRegistry).datedKeyValueStore.(*mongoDatedKeyValueStore).DB("test_driver_mongo").DropDatabase()

	if _, err := reg.(*datedServiceKeyRegistry).stampedPut("a_b-c", "kore ga kagi dayo."); err != nil {
		t.Fatal(err)
	}

	testDatedServiceKeyRegistry(t, reg)
}
