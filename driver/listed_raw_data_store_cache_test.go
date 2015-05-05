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
	"reflect"
	"testing"
	"time"
)

func TestCachingListedRawDataStore(t *testing.T) {
	testListedRawDataStore(t, newCachingListedRawDataStore(newMemoryListedRawDataStore(0, 0)))
}

func TestCachingListedRawDataStoreStamp(t *testing.T) {
	testListedRawDataStoreStamp(t, newCachingListedRawDataStore(newMemoryListedRawDataStore(0, 0)))
}

func TestCachingListedRawDataStoreExpiration(t *testing.T) {
	staleDur := 10 * time.Millisecond
	expiDur := 50 * time.Millisecond
	drv := newCachingListedRawDataStore(newMemoryListedRawDataStore(staleDur, expiDur))
	defer drv.Close()

	// 入れる。
	if _, err := drv.Put(testKey, testData); err != nil {
		t.Fatal(err)
	}
	if _, err := drv.Put(testKey+"a", testData); err != nil {
		t.Fatal(err)
	}

	end := time.Now().Add(2 * expiDur)
	var caData []byte
	var caDataStmp *Stamp
	var caKeys map[string]bool
	var caKeysStmp *Stamp
	for time.Now().Before(end) {
		data, newCaStmp, err := drv.Get(testKey, caDataStmp)
		if err != nil {
			t.Fatal(err)
		}
		if data != nil {
			caData = data
		}
		caDataStmp = newCaStmp
		if !reflect.DeepEqual(caData, testData) {
			t.Fatal(caData)
		}

		keys, newCaStmp, err := drv.Keys(caKeysStmp)
		if err != nil {
			t.Fatal(err)
		}
		if keys != nil {
			caKeys = keys
		}
		caKeysStmp = newCaStmp
		if len(caKeys) != 2 {
			t.Fatal(caKeys)
		}

		time.Sleep(staleDur / 3)
	}
}
