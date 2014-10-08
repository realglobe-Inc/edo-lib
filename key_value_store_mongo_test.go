package driver

import (
	"testing"
)

func TestMongoKeyValueStore(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := newMongoKeyValueStore(mongoAddr, testLabel, "key-value-store", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.base.DB(testLabel).DropDatabase()

	testKeyValueStore(t, reg)
}

func TestMongoKeyValueStoreStamp(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := newMongoKeyValueStore(mongoAddr, testLabel, "key-value-store", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.base.DB(testLabel).DropDatabase()

	testKeyValueStoreStamp(t, reg)
}
