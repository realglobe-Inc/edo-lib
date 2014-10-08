package driver

import (
	"reflect"
	"testing"
	"time"
)

func testTimeLimitedKeyValueStore(t *testing.T, reg TimeLimitedKeyValueStore) {
	expiDur := 10 * time.Millisecond

	// まだ無い。
	value1, _, err := reg.Get(testKey, nil)
	if err != nil {
		t.Fatal(err)
	} else if value1 != nil {
		t.Error(value1)
	}

	// 入れる。
	if _, err := reg.Put(testKey, testValue, time.Now().Add(expiDur)); err != nil {
		t.Fatal(err)
	}

	// ある。
	value2, _, err := reg.Get(testKey, nil)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(value2, testValue) {
		if !jsonEqual(value2, testValue) {
			t.Error(value2)
		}
	}

	// 消す。
	if err := reg.Remove(testKey); err != nil {
		t.Fatal(err)
	}

	// もう無い。
	value3, _, err := reg.Get(testKey, nil)
	if err != nil {
		t.Fatal(err)
	} else if value3 != nil {
		t.Error(value3)
	}

	// また入れる。
	if _, err := reg.Put(testKey, testValue, time.Now().Add(expiDur)); err != nil {
		t.Fatal(err)
	}

	// ある。
	value4, _, err := reg.Get(testKey, nil)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(value4, testValue) {
		if !jsonEqual(value4, testValue) {
			t.Error(value4)
		}
	}

	// 消えるまで待つ。
	time.Sleep(2 * expiDur)

	// もう無い。
	value5, _, err := reg.Get(testKey, nil)
	if err != nil {
		t.Fatal(err)
	} else if value5 != nil {
		t.Error(value5)
	}
}
