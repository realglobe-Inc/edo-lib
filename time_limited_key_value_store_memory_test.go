package driver

import (
	"testing"
)

func TestMemoryTimeLimitedKeyValueStore(t *testing.T) {
	testTimeLimitedKeyValueStore(t, NewMemoryTimeLimitedKeyValueStore())
}
