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

	reg := NewMongoKeyValueStore(mongoAddr, testLabel, "key-value-store", "key", nil, func(val interface{}) (interface{}, error) {
		delete(val.(map[string]interface{}), "_id")
		return val, nil
	}, nil, nil, 0, 0)
	defer reg.Clear()

	// まだ無い。
	val1, _, err := reg.Get(testKey, nil)
	if err != nil {
		t.Fatal(err)
	} else if val1 != nil {
		t.Error(val1)
	}

	// 入れる。
	now := time.Now() // mongodb の時間粒度がミリ秒なので細工する。
	val := map[string]interface{}{
		"key":    testKey,
		"date":   time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond()-now.Nanosecond()%1000000, now.Location()),
		"digest": "abcde",
		"array":  []interface{}{"elem-1", "elem-2"},
	}
	if _, err := reg.Put(testKey, val); err != nil {
		t.Fatal(err)
	}

	// ある。
	val2, _, err := reg.Get(testKey, nil)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(val2, val) {
		if !jsonEqual(val2, val) {
			t.Error(val2, val)
		}
	}

	// 消す。
	if err := reg.Remove(testKey); err != nil {
		t.Fatal(err)
	}

	// もう無い。
	val3, _, err := reg.Get(testKey, nil)
	if err != nil {
		t.Fatal(err)
	} else if val3 != nil {
		t.Error(val3)
	}
}

func TestMongoKeyValueStoreStamp(t *testing.T) {
	// ////////////////////////////////
	// util.SetupConsoleLog("github.com/realglobe-Inc", level.ALL)
	// defer util.SetupConsoleLog("github.com/realglobe-Inc", level.OFF)
	// ////////////////////////////////
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg := NewMongoKeyValueStore(mongoAddr, testLabel, "key-value-store", "key", nil, func(val interface{}) (interface{}, error) {
		delete(val.(map[string]interface{}), "_id")
		return val, nil
	}, nil, nil, 0, 0)
	defer reg.Clear()

	// まだ無い。
	val1, stmp1, err := reg.Get(testKey, nil)
	if err != nil {
		t.Fatal(err)
	} else if val1 != nil || stmp1 != nil {
		t.Error(val1, stmp1)
	}

	// 入れる。
	now := time.Now() // mongodb の時間粒度がミリ秒なので細工する。
	val := map[string]interface{}{
		"key":    testKey,
		"date":   time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond()-now.Nanosecond()%1000000, now.Location()),
		"digest": "abcde",
		"array":  []interface{}{"elem-1", "elem-2"},
	}
	stmp2, err := reg.Put(testKey, val)
	if err != nil {
		t.Fatal(err)
	}

	// ある。
	val3, stmp3, err := reg.Get(testKey, nil)
	if err != nil {
		t.Fatal(err)
	} else if stmp3 == nil {
		t.Error(stmp3)
	} else if !reflect.DeepEqual(val3, val) {
		if !jsonEqual(val3, val) {
			t.Error(val3)
		}
	}

	// キャッシュと同じだから返らない。
	val4, stmp4, err := reg.Get(testKey, stmp2)
	if err != nil {
		t.Fatal(err)
	} else if val4 != nil || stmp4 == nil {
		t.Error(val4, stmp4)
	}

	// キャッシュが古いから返る。
	val5, stmp5, err := reg.Get(testKey, &Stamp{Date: stmp2.Date.Add(-time.Second), Digest: stmp2.Digest})
	if err != nil {
		t.Fatal(err)
	} else if stmp5 == nil {
		t.Error(stmp5)
	} else if !reflect.DeepEqual(val5, val) {
		if !jsonEqual(val5, val) {
			t.Error(val5)
		}
	}

	// ダイジェストが違うから返る。
	val6, stmp6, err := reg.Get(testKey, &Stamp{Date: stmp2.Date, Digest: stmp2.Digest + "a"})
	if err != nil {
		t.Fatal(err)
	} else if stmp6 == nil {
		t.Error(stmp6)
	} else if !reflect.DeepEqual(val6, val) {
		if !jsonEqual(val6, val) {
			t.Error(val6)
		}
	}

	// 消す。
	if err := reg.Remove(testKey); err != nil {
		t.Fatal(err)
	}

	// もう無い。
	val7, stmp7, err := reg.Get(testKey, stmp2)
	if err != nil {
		t.Fatal(err)
	} else if val7 != nil || stmp7 != nil {
		t.Error(val7, stmp7)
	}
}

func TestMongoNKeyValueStore(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg := NewMongoNKeyValueStore(mongoAddr, testLabel, "key-value-store", []string{"key1", "key2"}, nil, func(val interface{}) (interface{}, error) {
		delete(val.(map[string]interface{}), "_id")
		return val, nil
	}, nil, nil, 0, 0)
	defer reg.Clear()

	testKey2 := testKey + "2"
	tagKeys := bson.M{"key1": testKey, "key2": testKey2}

	// まだ無い。
	val1, _, err := reg.NGet(tagKeys, nil)
	if err != nil {
		t.Fatal(err)
	} else if val1 != nil {
		t.Error(val1)
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
	if _, err := reg.NPut(tagKeys, val); err != nil {
		t.Fatal(err)
	}

	// ある。
	val2, _, err := reg.NGet(tagKeys, nil)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(val2, val) {
		if !jsonEqual(val2, val) {
			t.Error(val2, val)
		}
	}
	// キーが 1 つ違うので無い。
	val3, _, err := reg.NGet(bson.M{"key1": testKey, "key2": testKey}, nil)
	if err != nil {
		t.Fatal(err)
	} else if val3 != nil {
		t.Error(val3, val)
	}

	// 消す。
	if err := reg.NRemove(tagKeys); err != nil {
		t.Fatal(err)
	}

	// もう無い。
	val4, _, err := reg.NGet(tagKeys, nil)
	if err != nil {
		t.Fatal(err)
	} else if val4 != nil {
		t.Error(val4)
	}
}
