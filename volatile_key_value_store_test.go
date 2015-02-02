package driver

import (
	"reflect"
	"testing"
	"time"
)

func testVolatileKeyValueStore(t *testing.T, reg VolatileKeyValueStore) {
	expiDur := 10 * time.Millisecond

	// まだ無い。
	if v, _, err := reg.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if v != nil {
		t.Error(v)
	}

	// 入れる。
	if _, err := reg.Put(testKey, testVal, time.Now().Add(expiDur)); err != nil {
		t.Fatal(err)
	}

	// ある。
	if v, _, err := reg.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(v, testVal) {
		if !jsonEqual(v, testVal) {
			t.Error(v)
		}
	}

	// 消す。
	if err := reg.Remove(testKey); err != nil {
		t.Fatal(err)
	}

	// もう無い。
	if v, _, err := reg.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if v != nil {
		t.Error(v)
	}

	// また入れる。
	if _, err := reg.Put(testKey, testVal, time.Now().Add(expiDur)); err != nil {
		t.Fatal(err)
	}

	// ある。
	if v, _, err := reg.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(v, testVal) {
		if !jsonEqual(v, testVal) {
			t.Error(v)
		}
	}

	// 消えるまで待つ。
	time.Sleep(2 * expiDur)

	// もう無い。
	if v, _, err := reg.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if v != nil {
		t.Error(v)
	}
}
