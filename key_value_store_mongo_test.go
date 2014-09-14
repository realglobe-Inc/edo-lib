package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMongoKeyValueStore(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := newMongoKeyValueStore(mongoAddr, testLabel, "key-value-store", "key", "value")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.DB(testLabel).DropDatabase()

	testKeyValueStore(t, reg)
}

// キャッシュ用。
func TestMongoDatedKeyValueStore(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := newMongoDatedKeyValueStore(mongoAddr, testLabel, "key-value-store", 0, "key", "value")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.DB(testLabel).DropDatabase()

	testDatedKeyValueStore(t, reg)
}
