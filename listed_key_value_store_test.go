package driver

import (
	"reflect"
	"testing"
	"time"
)

func testListedKeyValueStore(t *testing.T, reg ListedKeyValueStore) {
	// まだ無い。
	val1, _, err := reg.Get(testKey, nil)
	if err != nil {
		t.Fatal(err)
	} else if val1 != nil {
		t.Error(val1)
	}

	// 入れる。
	if _, err := reg.Put(testKey, testVal); err != nil {
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
	val3, _, err := reg.Get(testKey, nil)
	if err != nil {
		t.Fatal(err)
	} else if val3 != nil {
		t.Error(val3)
	}
}

func testListedKeyValueStoreStamp(t *testing.T, reg ListedKeyValueStore) {
	// まだ無い。
	val1, stmp1, err := reg.Get(testKey, nil)
	if err != nil {
		t.Fatal(err)
	} else if val1 != nil || stmp1 != nil {
		t.Error(val1, stmp1)
	}

	// 入れる。
	stmp2, err := reg.Put(testKey, testVal)
	if err != nil {
		t.Fatal(err)
	}

	// ある。
	val3, stmp3, err := reg.Get(testKey, nil)
	if err != nil {
		t.Fatal(err)
	} else if stmp3 == nil {
		t.Error(stmp3)
	} else if !reflect.DeepEqual(val3, testVal) {
		if !jsonEqual(val3, testVal) {
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
	} else if !reflect.DeepEqual(val5, testVal) {
		if !jsonEqual(val5, testVal) {
			t.Error(val5)
		}
	}

	// ダイジェストが違うから返る。
	val6, stmp6, err := reg.Get(testKey, &Stamp{Date: stmp2.Date, Digest: stmp2.Digest + "a"})
	if err != nil {
		t.Fatal(err)
	} else if stmp6 == nil {
		t.Error(stmp6)
	} else if !reflect.DeepEqual(val6, testVal) {
		if !jsonEqual(val6, testVal) {
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
