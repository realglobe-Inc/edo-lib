package driver

import (
	"reflect"
	"testing"
	"time"
)

func testKeyValueStore(t *testing.T, reg KeyValueStore) {
	// まだ無い。
	value1, _, err := reg.Get(testKey, nil)
	if err != nil {
		t.Fatal(err)
	} else if value1 != nil {
		t.Error(value1)
	}

	// 入れる。
	if _, err := reg.Put(testKey, testValue); err != nil {
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
}

func testKeyValueStoreStamp(t *testing.T, reg KeyValueStore) {
	// まだ無い。
	value1, stmp1, err := reg.Get(testKey, nil)
	if err != nil {
		t.Fatal(err)
	} else if value1 != nil || stmp1 != nil {
		t.Error(value1, stmp1)
	}

	// 入れる。
	stmp2, err := reg.Put(testKey, testValue)
	if err != nil {
		t.Fatal(err)
	}

	// ある。
	value3, stmp3, err := reg.Get(testKey, nil)
	if err != nil {
		t.Fatal(err)
	} else if stmp3 == nil {
		t.Error(stmp3)
	} else if !reflect.DeepEqual(value3, testValue) {
		if !jsonEqual(value3, testValue) {
			t.Error(value3)
		}
	}

	// キャッシュと同じだから返らない。
	value4, stmp4, err := reg.Get(testKey, stmp2)
	if err != nil {
		t.Fatal(err)
	} else if value4 != nil || stmp4 == nil {
		t.Error(value4, stmp4)
	}

	// キャッシュが古いから返る。
	value5, stmp5, err := reg.Get(testKey, &Stamp{Date: stmp2.Date.Add(-time.Second), Digest: stmp2.Digest})
	if err != nil {
		t.Fatal(err)
	} else if stmp5 == nil {
		t.Error(stmp5)
	} else if !reflect.DeepEqual(value5, testValue) {
		if !jsonEqual(value5, testValue) {
			t.Error(value5)
		}
	}

	// ダイジェストが違うから返る。
	value6, stmp6, err := reg.Get(testKey, &Stamp{Date: stmp2.Date, Digest: stmp2.Digest + "a"})
	if err != nil {
		t.Fatal(err)
	} else if stmp6 == nil {
		t.Error(stmp6)
	} else if !reflect.DeepEqual(value6, testValue) {
		if !jsonEqual(value6, testValue) {
			t.Error(value6)
		}
	}

	// 消す。
	if err := reg.Remove(testKey); err != nil {
		t.Fatal(err)
	}

	// もう無い。
	value7, stmp7, err := reg.Get(testKey, stmp2)
	if err != nil {
		t.Fatal(err)
	} else if value7 != nil || stmp7 != nil {
		t.Error(value7, stmp7)
	}
}
