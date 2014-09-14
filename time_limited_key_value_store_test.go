package driver

import (
	"reflect"
	"testing"
	"time"
)

// 非キャッシュ用。
func testTimeLimitedKeyValueStore(t *testing.T, reg TimeLimitedKeyValueStore) {
	expiDur := 50 * time.Millisecond

	// まだ無い。
	value1, err := reg.Get(testKey)
	if err != nil {
		t.Fatal(err)
	} else if value1 != nil {
		t.Error(value1)
	}

	// 入れる。
	if err := reg.Put(testKey, testValue, time.Now().Add(expiDur)); err != nil {
		t.Fatal(err)
	}

	// ある。
	value2, err := reg.Get(testKey)
	if err != nil {
		t.Fatal(err)
	} else if value2 == nil || !reflect.DeepEqual(value2, testValue) {
		t.Error(value2)
	}

	// 消す。
	if err := reg.Remove(testKey); err != nil {
		t.Fatal(err)
	}

	// もう無い。
	value3, err := reg.Get(testKey)
	if err != nil {
		t.Fatal(err)
	} else if value3 != nil {
		t.Error(value3)
	}

	// また入れる。
	if err := reg.Put(testKey, testValue, time.Now().Add(expiDur)); err != nil {
		t.Fatal(err)
	}

	// ある。
	value4, err := reg.Get(testKey)
	if err != nil {
		t.Fatal(err)
	} else if value4 == nil || !reflect.DeepEqual(value4, testValue) {
		t.Error(value4)
	}

	// 消えるまで待つ。
	time.Sleep(2 * expiDur)

	// もう無い。
	value5, err := reg.Get(testKey)
	if err != nil {
		t.Fatal(err)
	} else if value5 != nil {
		t.Error(value5)
	}
}
