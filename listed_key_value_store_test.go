package driver

import (
	"reflect"
	"testing"
	"time"
)

func testListedKeyValueStore(t *testing.T, reg ListedKeyValueStore) {
	// まだ無い。
	if v, _, err := reg.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if v != nil {
		t.Error(v)
	}

	// 入れる。
	if _, err := reg.Put(testKey, testVal); err != nil {
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

	keys, _, err := reg.Keys(nil)
	if err != nil {
		t.Fatal(err)
	} else if len(keys) != 1 || !keys[testKey] {
		t.Error(keys)
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
}

func testListedKeyValueStoreStamp(t *testing.T, reg ListedKeyValueStore) {
	// まだ無い。
	if v, s, err := reg.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if v != nil || s != nil {
		t.Error(v, s)
	}

	// 入れる。
	stmp, err := reg.Put(testKey, testVal)
	if err != nil {
		t.Fatal(err)
	}

	// ある。
	if v, s, err := reg.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if s == nil {
		t.Error(s)
	} else if !reflect.DeepEqual(v, testVal) {
		if !jsonEqual(v, testVal) {
			t.Error(v)
		}
	}

	// キャッシュと同じだから返らない。
	if v, s, err := reg.Get(testKey, stmp); err != nil {
		t.Fatal(err)
	} else if v != nil || s == nil {
		t.Error(v, s)
	}

	// キャッシュが古いから返る。
	if v, s, err := reg.Get(testKey, &Stamp{Date: stmp.Date.Add(-time.Second), Digest: stmp.Digest}); err != nil {
		t.Fatal(err)
	} else if s == nil {
		t.Error(s)
	} else if !reflect.DeepEqual(v, testVal) {
		if !jsonEqual(v, testVal) {
			t.Error(v)
		}
	}

	// ダイジェストが違うから返る。
	if v, s, err := reg.Get(testKey, &Stamp{Date: stmp.Date, Digest: stmp.Digest + "a"}); err != nil {
		t.Fatal(err)
	} else if s == nil {
		t.Error(s)
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
	if v, s, err := reg.Get(testKey, stmp); err != nil {
		t.Fatal(err)
	} else if v != nil || s != nil {
		t.Error(v, s)
	}
}
