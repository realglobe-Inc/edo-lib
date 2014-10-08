package driver

import (
	"testing"
)

func TestWebTimeLimitedKeyValueStore(t *testing.T) {
	testTimeLimitedKeyValueStore(t, newSynchronizedTimeLimitedKeyValueStore(newMemoryTimeLimitedKeyValueStore(0)))
}
