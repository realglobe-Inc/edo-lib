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
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"testing"
	"time"
)

func TestMongoKeyValueStore(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	drv := NewMongoKeyValueStore(mongoAddr, testLabel, "key-value-store", "key", nil, func(val interface{}) (interface{}, error) {
		delete(val.(map[string]interface{}), "_id")
		return val, nil
	}, nil, nil, 0, 0)
	defer drv.Close()
	defer drv.Clear()

	// まだ無い。
	if v, _, err := drv.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if v != nil {
		t.Error(v)
	}

	// 入れる。
	now := time.Now() // mongodb の時間粒度がミリ秒なので細工する。
	val := map[string]interface{}{
		"key":    testKey,
		"date":   time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond()-now.Nanosecond()%1000000, now.Location()),
		"digest": "abcde",
		"array":  []interface{}{"elem-1", "elem-2"},
	}
	if _, err := drv.Put(testKey, val); err != nil {
		t.Fatal(err)
	}

	// ある。
	if v, _, err := drv.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(v, val) {
		if !jsonEqual(v, val) {
			t.Error(v, val)
		}
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

func TestMongoKeyValueStoreStamp(t *testing.T) {
	// ////////////////////////////////
	// logutil.SetupConsole("github.com/realglobe-Inc", level.ALL)
	// defer logutil.SetupConsole("github.com/realglobe-Inc", level.OFF)
	// ////////////////////////////////
	if mongoAddr == "" {
		t.SkipNow()
	}

	drv := NewMongoKeyValueStore(mongoAddr, testLabel, "key-value-store", "key", nil, func(val interface{}) (interface{}, error) {
		delete(val.(map[string]interface{}), "_id")
		return val, nil
	}, nil, nil, 0, 0)
	defer drv.Close()
	defer drv.Clear()

	// まだ無い。
	if v, s, err := drv.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if v != nil || s != nil {
		t.Error(v, s)
	}

	// 入れる。
	now := time.Now() // mongodb の時間粒度がミリ秒なので細工する。
	val := map[string]interface{}{
		"key":    testKey,
		"date":   time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond()-now.Nanosecond()%1000000, now.Location()),
		"digest": "abcde",
		"array":  []interface{}{"elem-1", "elem-2"},
	}
	stmp, err := drv.Put(testKey, val)
	if err != nil {
		t.Fatal(err)
	}

	// ある。
	if v, s, err := drv.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if s == nil {
		t.Error(s)
	} else if !reflect.DeepEqual(v, val) {
		if !jsonEqual(v, val) {
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
	} else if !reflect.DeepEqual(v, val) {
		if !jsonEqual(v, val) {
			t.Error(v)
		}
	}

	// ダイジェストが違うから返る。
	if v, s, err := drv.Get(testKey, &Stamp{Date: stmp.Date, Digest: stmp.Digest + "a"}); err != nil {
		t.Fatal(err)
	} else if s == nil {
		t.Error(s)
	} else if !reflect.DeepEqual(v, val) {
		if !jsonEqual(v, val) {
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

func TestMongoNKeyValueStore(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	drv := NewMongoNKeyValueStore(mongoAddr, testLabel, "key-value-store", []string{"key1", "key2"}, nil, func(val interface{}) (interface{}, error) {
		delete(val.(map[string]interface{}), "_id")
		return val, nil
	}, nil, nil, 0, 0)
	defer drv.Close()
	defer drv.Clear()

	testKey2 := testKey + "2"
	tagKeys := bson.M{"key1": testKey, "key2": testKey2}

	// まだ無い。
	if v, _, err := drv.NGet(tagKeys, nil); err != nil {
		t.Fatal(err)
	} else if v != nil {
		t.Error(v)
	}

	// 入れる。
	now := time.Now() // mongodb の時間粒度がミリ秒なので細工する。
	val := map[string]interface{}{
		"key1":   testKey,
		"key2":   testKey2,
		"date":   time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond()-now.Nanosecond()%1000000, now.Location()),
		"digest": "abcde",
		"array":  []interface{}{"elem-1", "elem-2"},
	}
	if _, err := drv.NPut(tagKeys, val); err != nil {
		t.Fatal(err)
	}

	// ある。
	if v, _, err := drv.NGet(tagKeys, nil); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(v, val) {
		if !jsonEqual(v, val) {
			t.Error(v, val)
		}
	}
	// キーが 1 つ違うので無い。
	if v, _, err := drv.NGet(bson.M{"key1": testKey, "key2": testKey}, nil); err != nil {
		t.Fatal(err)
	} else if v != nil {
		t.Error(v, val)
	}

	// 消す。
	if err := drv.NRemove(tagKeys); err != nil {
		t.Fatal(err)
	}

	// もう無い。
	if v, _, err := drv.NGet(tagKeys, nil); err != nil {
		t.Fatal(err)
	} else if v != nil {
		t.Error(v)
	}
}
