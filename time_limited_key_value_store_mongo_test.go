package driver

import (
	"testing"
)

func TestMongoTimeLimitedKeyValueStore(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg := newMongoTimeLimitedKeyValueStore(mongoAddr, testLabel, "time-limited-key-value-store", nil, nil, nil, 0, 0)
	defer reg.Clear()

	testTimeLimitedKeyValueStore(t, reg)
}
