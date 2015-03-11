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

func testListedKeyValueStore(t *testing.T, drv ListedKeyValueStore) {
	defer drv.Close()

	// まだ無い。
	if v, _, err := drv.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if v != nil {
		t.Error(v)
	}

	// 入れる。
	if _, err := drv.Put(testKey, testVal); err != nil {
		t.Fatal(err)
	}

	// ある。
	if v, _, err := drv.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(v, testVal) {
		if !jsonEqual(v, testVal) {
			t.Error(v)
		}
	}

	keys, _, err := drv.Keys(nil)
	if err != nil {
		t.Fatal(err)
	} else if len(keys) != 1 || !keys[testKey] {
		t.Error(keys)
	}

	// 消す。
	if err := drv.Remove(testKey); err != nil {
		t.Fatal(err)
	}

	// もう無い。
	if v, _, err := drv.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if v != nil {
		t.Error(v)
	}
}

func testListedKeyValueStoreStamp(t *testing.T, drv ListedKeyValueStore) {
	defer drv.Close()

	// まだ無い。
	if v, s, err := drv.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if v != nil || s != nil {
		t.Error(v, s)
	}

	// 入れる。
	stmp, err := drv.Put(testKey, testVal)
	if err != nil {
		t.Fatal(err)
	}

	// ある。
	if v, s, err := drv.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if s == nil {
		t.Error(s)
	} else if !reflect.DeepEqual(v, testVal) {
		if !jsonEqual(v, testVal) {
			t.Error(v)
		}
	}

	// キャッシュと同じだから返らない。
	if v, s, err := drv.Get(testKey, stmp); err != nil {
		t.Fatal(err)
	} else if v != nil || s == nil {
		t.Error(v, s)
	}

	// キャッシュが古いから返る。
	if v, s, err := drv.Get(testKey, &Stamp{Date: stmp.Date.Add(-time.Second), Digest: stmp.Digest}); err != nil {
		t.Fatal(err)
	} else if s == nil {
		t.Error(s)
	} else if !reflect.DeepEqual(v, testVal) {
		if !jsonEqual(v, testVal) {
			t.Error(v)
		}
	}

	// ダイジェストが違うから返る。
	if v, s, err := drv.Get(testKey, &Stamp{Date: stmp.Date, Digest: stmp.Digest + "a"}); err != nil {
		t.Fatal(err)
	} else if s == nil {
		t.Error(s)
	} else if !reflect.DeepEqual(v, testVal) {
		if !jsonEqual(v, testVal) {
			t.Error(v)
		}
	}

	// 消す。
	if err := drv.Remove(testKey); err != nil {
		t.Fatal(err)
	}

	// もう無い。
	if v, s, err := drv.Get(testKey, stmp); err != nil {
		t.Fatal(err)
	} else if v != nil || s != nil {
		t.Error(v, s)
	}
}
