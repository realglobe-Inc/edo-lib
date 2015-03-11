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
	"time"
)

type memoryListedRawDataStore memoryListedKeyValueStore

// スレッドセーフ。
func NewMemoryListedRawDataStore(staleDur, expiDur time.Duration) ListedRawDataStore {
	return newSynchronizedListedRawDataStore(newMemoryListedRawDataStore(staleDur, expiDur))
}

// スレッドセーフではない。
func newMemoryListedRawDataStore(staleDur, expiDur time.Duration) *memoryListedRawDataStore {
	return (*memoryListedRawDataStore)(newMemoryListedKeyValueStore(staleDur, expiDur))
}

func (drv *memoryListedRawDataStore) Keys(caStmp *Stamp) (keys map[string]bool, newCaStmp *Stamp, err error) {
	return ((*memoryListedKeyValueStore)(drv)).Keys(caStmp)
}

func (drv *memoryListedRawDataStore) Get(key string, caStmp *Stamp) (data []byte, newCaStmp *Stamp, err error) {
	val, newCaStmp, err := ((*memoryListedKeyValueStore)(drv)).Get(key, caStmp)
	if val == nil {
		return nil, newCaStmp, err
	}
	return val.([]byte), newCaStmp, nil
}

func (drv *memoryListedRawDataStore) Put(key string, data []byte) (*Stamp, error) {
	return ((*memoryListedKeyValueStore)(drv)).Put(key, data)
}

func (drv *memoryListedRawDataStore) Remove(key string) error {
	return ((*memoryListedKeyValueStore)(drv)).Remove(key)
}

func (drv *memoryListedRawDataStore) Close() error {
	return ((*memoryListedKeyValueStore)(drv)).Close()
}
