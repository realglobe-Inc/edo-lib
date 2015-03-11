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

func testRawDataStore(t *testing.T, drv RawDataStore) {
	defer drv.Close()

	// まだ無い。
	if d, _, err := drv.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if d != nil {
		t.Error(d)
	}

	// 入れる。
	if _, err := drv.Put(testKey, testData); err != nil {
		t.Fatal(err)
	}

	// ある。
	if d, _, err := drv.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if d == nil || !reflect.DeepEqual(d, testData) {
		t.Error(d)
	}

	// 消す。
	if err := drv.Remove(testKey); err != nil {
		t.Fatal(err)
	}

	// もう無い。
	if d, _, err := drv.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if d != nil {
		t.Error(d)
	}
}

func testRawDataStoreStamp(t *testing.T, drv RawDataStore) {
	defer drv.Close()

	// まだ無い。
	if d, s, err := drv.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if d != nil || s != nil {
		t.Error(d, s)
	}

	// 入れる。
	stmp, err := drv.Put(testKey, testData)
	if err != nil {
		t.Fatal(err)
	}

	// ある。
	if d, s, err := drv.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if d == nil || !reflect.DeepEqual(d, testData) || s == nil {
		t.Error(d, s)
	}

	// キャッシュと同じだから返らない。
	if d, s, err := drv.Get(testKey, stmp); err != nil {
		t.Fatal(err)
	} else if d != nil || s == nil {
		t.Error(d, s)
	}

	// キャッシュが古いから返る。
	if d, s, err := drv.Get(testKey, &Stamp{Date: stmp.Date.Add(-time.Second), Digest: stmp.Digest}); err != nil {
		t.Fatal(err)
	} else if d == nil || !reflect.DeepEqual(d, testData) || s == nil {
		t.Error(d, s)
	}

	// ダイジェストが違うから返る。
	if d, s, err := drv.Get(testKey, &Stamp{Date: stmp.Date, Digest: stmp.Digest + "a"}); err != nil {
		t.Fatal(err)
	} else if d == nil || !reflect.DeepEqual(d, testData) || s == nil {
		t.Error(d, s)
	}

	// 消す。
	if err := drv.Remove(testKey); err != nil {
		t.Fatal(err)
	}

	// もう無い。
	if d, s, err := drv.Get(testKey, stmp); err != nil {
		t.Fatal(err)
	} else if d != nil || s != nil {
		t.Error(d, s)
	}
}
