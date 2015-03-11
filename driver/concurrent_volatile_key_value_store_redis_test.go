// Copyright 2015 realglobe, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
