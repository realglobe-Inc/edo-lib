package driver

import (
	"reflect"
	"testing"
	"time"
)

func testListedRawDataStore(t *testing.T, drv ListedRawDataStore) {
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
	if d, _, err := drv.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if d != nil {
		t.Error(d)
	}
}

func testListedRawDataStoreStamp(t *testing.T, drv ListedRawDataStore) {
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
