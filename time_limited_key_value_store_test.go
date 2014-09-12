package driver

import (
	"testing"
	"time"
)

func testTimeLimitedKeyValueStore(t *testing.T, reg TimeLimitedKeyValueStore) {
	key := "abcdAbcd1234-+/:"
	value := "aaa"
	expiDur := 100 * time.Millisecond

	// まだ無い。
	value1, err := reg.Get(key)
	if err != nil {
		t.Fatal(err)
	} else if value1 != nil {
		t.Error(value1)
	}

	// 入れる。
	if err := reg.Put(key, value, time.Now().Add(expiDur)); err != nil {
		t.Fatal(err)
	}

	// ある。
	value2, err := reg.Get(key)
	if err != nil {
		t.Fatal(err)
	} else if value2 == nil || value2.(string) != value {
		t.Error(value2)
	}

	// 消す。
	if err := reg.Remove(key); err != nil {
		t.Fatal(err)
	}

	// もう無い。
	value3, err := reg.Get(key)
	if err != nil {
		t.Fatal(err)
	} else if value3 != nil {
		t.Error(value3)
	}

	// また入れる。
	if err := reg.Put(key, value, time.Now().Add(expiDur)); err != nil {
		t.Fatal(err)
	}

	// ある。
	value4, err := reg.Get(key)
	if err != nil {
		t.Fatal(err)
	} else if value4 == nil || value4.(string) != value {
		t.Error(value4)
	}

	// 消えるまで待つ。
	time.Sleep(2 * expiDur)

	// もう無い。
	value5, err := reg.Get(key)
	if err != nil {
		t.Fatal(err)
	} else if value5 != nil {
		t.Error(value5)
	}
}
