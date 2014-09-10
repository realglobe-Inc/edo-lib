package driver

import (
	"testing"
	"time"
)

// 非キャッシュ用。
func testKeyValueStore(t *testing.T, reg keyValueStore) {
	key := "abcdAbcd1234-+/:"
	value := "aaa"

	// まだ無い。
	value1, err := reg.get(key)
	if err != nil {
		t.Fatal(err)
	} else if value1 != nil {
		t.Error(value1)
	}

	// 入れる。
	if err := reg.put(key, value); err != nil {
		t.Fatal(err)
	}

	// ある。
	value2, err := reg.get(key)
	if err != nil {
		t.Fatal(err)
	} else if value2 == nil || value2.(string) != value {
		t.Error(value2)
	}

	// 消す。
	if err := reg.remove(key); err != nil {
		t.Fatal(err)
	}

	// もう無い。
	value3, err := reg.get(key)
	if err != nil {
		t.Fatal(err)
	} else if value3 != nil {
		t.Error(value3)
	}
}

// キャッシュ用。
func testDatedKeyValueStore(t *testing.T, reg datedKeyValueStore) {
	key := "abcdAbcd1234-+/:"
	value := "aaa"

	// まだ無い。
	value1, stmp1, err := reg.stampedGet(key, nil)
	if err != nil {
		t.Fatal(err)
	} else if (value1 != nil && value1 != "") || stmp1 != nil {
		t.Error(value1, stmp1)
	}

	// 入れる。
	stmp2, err := reg.stampedPut(key, value)
	if err != nil {
		t.Fatal(err)
	}

	// ある。
	value3, stmp3, err := reg.stampedGet(key, nil)
	if err != nil {
		t.Fatal(err)
	} else if value3 == nil || value3.(string) != value || stmp3 == nil {
		t.Error(value3, stmp3)
	}

	// キャッシュと同じだから返らない。
	value4, stmp4, err := reg.stampedGet(key, stmp2)
	if err != nil {
		t.Fatal(err)
	} else if (value4 != nil && value4 != "") || stmp4 == nil {
		t.Error(value4, stmp4)
	}

	// キャッシュが古いから返る。
	value5, stmp5, err := reg.stampedGet(key, &Stamp{Date: stmp2.Date.Add(-time.Second), Digest: stmp2.Digest})
	if err != nil {
		t.Fatal(err)
	} else if value5 == nil || value5.(string) != value || stmp5 == nil {
		t.Error(value5, stmp5)
	}

	// ダイジェストが違うから返る。
	value6, stmp6, err := reg.stampedGet(key, &Stamp{Date: stmp2.Date, Digest: stmp2.Digest + "a"})
	if err != nil {
		t.Fatal(err)
	} else if value6 == nil || value6.(string) != value || stmp6 == nil {
		t.Error(value6, stmp6)
	}

	// 消す。
	if err := reg.remove(key); err != nil {
		t.Fatal(err)
	}

	// もう無い。
	value7, stmp7, err := reg.stampedGet(key, stmp2)
	if err != nil {
		t.Fatal(err)
	} else if (value7 != nil && value7 != "") || stmp7 != nil {
		t.Error(value7, stmp7)
	}
}
