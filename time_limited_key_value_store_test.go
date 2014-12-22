package driver

import (
	"reflect"
	"testing"
	"time"
)

func testTimeLimitedKeyValueStore(t *testing.T, reg TimeLimitedKeyValueStore) {
	expiDur := 10 * time.Millisecond

	// まだ無い。
	val1, _, err := reg.Get(testKey, nil)
	if err != nil {
		t.Fatal(err)
	} else if val1 != nil {
		t.Error(val1)
	}

	// 入れる。
	if _, err := reg.Put(testKey, testVal, time.Now().Add(expiDur)); err != nil {
		t.Fatal(err)
	}

	// ある。
	val2, _, err := reg.Get(testKey, nil)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(val2, testVal) {
		if !jsonEqual(val2, testVal) {
			t.Error(val2)
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

	// また入れる。
	if _, err := reg.Put(testKey, testVal, time.Now().Add(expiDur)); err != nil {
		t.Fatal(err)
	}

	// ある。
	val4, _, err := reg.Get(testKey, nil)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(val4, testVal) {
		if !jsonEqual(val4, testVal) {
			t.Error(val4)
		}
	}

	// 消えるまで待つ。
	time.Sleep(2 * expiDur)

	// もう無い。
	val5, _, err := reg.Get(testKey, nil)
	if err != nil {
		t.Fatal(err)
	} else if val5 != nil {
		t.Error(val5)
	}
}
