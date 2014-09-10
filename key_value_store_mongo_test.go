package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMongoKeyValueStore(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := newMongoKeyValueStore(mongoAddr, "test_driver_mongo", "kvs", "key", "value")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.DB("test_driver_mongo").DropDatabase()

	testKeyValueStore(t, reg)
}

// キャッシュ用。
func TestMongoDatedKeyValueStore(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := newMongoDatedKeyValueStore(mongoAddr, "test_driver_mongo", "kvs", 0, "key", "value")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.DB("test_driver_mongo").DropDatabase()

	testDatedKeyValueStore(t, reg)
}
