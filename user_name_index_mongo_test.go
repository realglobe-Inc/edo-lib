package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMongoUserNameIndex(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoUserNameIndex(mongoAddr, "test_driver_mongo", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*userNameIndex).keyValueStore.(*mongoKeyValueStore).DB("test_driver_mongo").DropDatabase()

	if err := reg.(*userNameIndex).put("a_b-c", "aaaa-bbbb-cccc"); err != nil {
		t.Fatal(err)
	}

	testUserNameIndex(t, reg)
}
