package driver

import (
	"testing"
)

func TestMongoTimeLimitedKeyValueStore(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := newMongoTimeLimitedKeyValueStore(mongoAddr, testLabel, "time-limited-key-value-store", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.base.base.DB(testLabel).DropDatabase()

	testTimeLimitedKeyValueStore(t, reg)
}
