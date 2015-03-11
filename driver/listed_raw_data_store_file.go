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

type fileListedRawDataStore struct {
	lister
	RawDataStore
}

// スレッドセーフ。
func NewFileListedRawDataStore(path string, keyToPath, pathToKey func(string) string, staleDur, expiDur time.Duration) ListedRawDataStore {
	return newSynchronizedListedRawDataStore(newCachingListedRawDataStore(newFileListedRawDataStore(path, keyToPath, pathToKey, staleDur, expiDur)))
}

// スレッドセーフではない。
func newFileListedRawDataStore(path string, keyToPath, pathToKey func(string) string, staleDur, expiDur time.Duration) *fileListedRawDataStore {
	return &fileListedRawDataStore{
		newFileLister(path, pathToKey, staleDur, expiDur),
		newFileRawDataStore(path, keyToPath, staleDur, expiDur),
	}
}
