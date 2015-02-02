package driver

import (
	"encoding/json"
	"testing"
	"time"
)

func TestRedisVolatileKeyValueStore(t *testing.T) {
	if redisAddr == "" {
		t.SkipNow()
	}

	testVolatileKeyValueStore(t, newRedisVolatileKeyValueStore(NewRedisPool(redisAddr, 5, time.Second), testLabel+":", json.Marshal, jsonUnmarshal, nil, time.Second, time.Second))
}
