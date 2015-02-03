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

	testVolatileKeyValueStore(t, newRedisConcurrentVolatileKeyValueStore(NewRedisPool(redisAddr, 5, time.Second), testLabel+":", json.Marshal, jsonUnmarshal, nil, time.Second, time.Second))
}

func TestRedisConcurrentVolatileKeyValueStore(t *testing.T) {
	if redisAddr == "" {
		t.SkipNow()
	}

	testConcurrentVolatileKeyValueStore(t, newRedisConcurrentVolatileKeyValueStore(NewRedisPool(redisAddr, 5, time.Second), testLabel+":", json.Marshal, jsonUnmarshal, nil, time.Second, time.Second))
}

func TestRedisConcurrentVolatileKeyValueStoreConsistency(t *testing.T) {
	if redisAddr == "" {
		t.SkipNow()
	}

	testConcurrentVolatileKeyValueStoreConsistency(t, NewRedisConcurrentVolatileKeyValueStore(NewRedisPool(redisAddr, 5, time.Second), testLabel+":", json.Marshal, jsonUnmarshal, nil, time.Second, time.Second))
}
