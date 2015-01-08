package driver

import (
	"encoding/json"
	"testing"
	"time"
)

func TestRedisTimeLimitedKeyValueStore(t *testing.T) {
	testTimeLimitedKeyValueStore(t, newRedisTimeLimitedKeyValueStore(NewRedisPool(redisAddr, 5, time.Second), testLabel+":", json.Marshal, jsonUnmarshal, nil, time.Second, time.Second))
}
