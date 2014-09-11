package driver

import (
	"testing"
)

func TestMongoUserAttributeRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoUserAttributeRegistry(mongoAddr, "test_driver_mongo", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*userAttributeRegistry).keyValueStore.(*mongoKeyValueStore).DB("test_driver_mongo").DropDatabase()

	if err := reg.(*userAttributeRegistry).put(userAttributeKey("a_b-c", "attribute"), "abcd"); err != nil {
		t.Fatal(err)
	}

	testUserAttributeRegistry(t, reg)
}
