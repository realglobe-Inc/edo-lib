package driver

import (
	"testing"
)

func TestMongoTimeLimitedKeyValueStore(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoTimeLimitedKeyValueStore(mongoAddr, "test_driver_mongo", "kvs", "key", "value")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoTimeLimitedKeyValueStore).DB("test_driver_mongo").DropDatabase()

	testTimeLimitedKeyValueStore(t, reg)
}
