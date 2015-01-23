package driver

import (
	"reflect"
	"testing"
	"time"
)

func testRawDataStore(t *testing.T, reg RawDataStore) {
	// まだ無い。
	if d, _, err := reg.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if d != nil {
		t.Error(d)
	}

	// 入れる。
	if _, err := reg.Put(testKey, testData); err != nil {
		t.Fatal(err)
	}

	// ある。
	if d, _, err := reg.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if d == nil || !reflect.DeepEqual(d, testData) {
		t.Error(d)
	}

	// 消す。
	if err := reg.Remove(testKey); err != nil {
		t.Fatal(err)
	}

	// もう無い。
	if d, _, err := reg.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if d != nil {
		t.Error(d)
	}
}

func testRawDataStoreStamp(t *testing.T, reg RawDataStore) {
	// まだ無い。
	if d, s, err := reg.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if d != nil || s != nil {
		t.Error(d, s)
	}

	// 入れる。
	stmp, err := reg.Put(testKey, testData)
	if err != nil {
		t.Fatal(err)
	}

	// ある。
	if d, s, err := reg.Get(testKey, nil); err != nil {
		t.Fatal(err)
	} else if d == nil || !reflect.DeepEqual(d, testData) || s == nil {
		t.Error(d, s)
	}

	// キャッシュと同じだから返らない。
	if d, s, err := reg.Get(testKey, stmp); err != nil {
		t.Fatal(err)
	} else if d != nil || s == nil {
		t.Error(d, s)
	}

	// キャッシュが古いから返る。
	if d, s, err := reg.Get(testKey, &Stamp{Date: stmp.Date.Add(-time.Second), Digest: stmp.Digest}); err != nil {
		t.Fatal(err)
	} else if d == nil || !reflect.DeepEqual(d, testData) || s == nil {
		t.Error(d, s)
	}

	// ダイジェストが違うから返る。
	if d, s, err := reg.Get(testKey, &Stamp{Date: stmp.Date, Digest: stmp.Digest + "a"}); err != nil {
		t.Fatal(err)
	} else if d == nil || !reflect.DeepEqual(d, testData) || s == nil {
		t.Error(d, s)
	}

	// 消す。
	if err := reg.Remove(testKey); err != nil {
		t.Fatal(err)
	}

	// もう無い。
	if d, s, err := reg.Get(testKey, stmp); err != nil {
		t.Fatal(err)
	} else if d != nil || s != nil {
		t.Error(d, s)
	}
}
