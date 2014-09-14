package driver

import (
	"testing"
)

// 非キャッシュ用。
func TestMongoTimeLimitedKeyValueStore(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoTimeLimitedKeyValueStore(mongoAddr, testLabel, "time-limited-key-value-store", "key", "value")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoTimeLimitedKeyValueStore).DB(testLabel).DropDatabase()

	testTimeLimitedKeyValueStore(t, reg)
}
