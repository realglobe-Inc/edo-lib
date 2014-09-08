package driver

import (
	"testing"
	"time"
)

// 事前に、UUID a_b-c、ユーザー属性取得用 URI https://localhost:1234/query で登録しとく。

// 非キャッシュ用。
func testServiceKeyRegistry(t *testing.T, reg ServiceKeyRegistry) {
	key, err := reg.ServiceKey("a_b-c")
	if err != nil {
		t.Fatal(err)
	} else if key != "kore ga kagi dayo." {
		t.Error(key)
	}
}

// キャッシュ用。
func testDatedServiceKeyRegistry(t *testing.T, reg DatedServiceKeyRegistry) {

	key1, stmp1, err := reg.StampedServiceKey("a_b-c", nil)
	if err != nil {
		t.Fatal(err)
	} else if key1 != "kore ga kagi dayo." || stmp1 == nil {
		t.Error(key1, stmp1)
	}

	// キャッシュと同じだから返らない。
	key2, stmp2, err := reg.StampedServiceKey("a_b-c", stmp1)
	if err != nil {
		t.Fatal(err)
	} else if key2 != "" || stmp2 == nil {
		t.Error(key2, stmp2)
	}

	// キャッシュが古いから返る。
	key3, stmp3, err := reg.StampedServiceKey("a_b-c", &Stamp{Date: stmp1.Date.Add(-time.Second), Digest: stmp1.Digest})
	if err != nil {
		t.Fatal(err)
	} else if key3 != "kore ga kagi dayo." || stmp3 == nil {
		t.Error(key3, stmp3)
	}

	// ダイジェストが違うから返る。
	key4, stmp4, err := reg.StampedServiceKey("a_b-c", &Stamp{Date: stmp1.Date, Digest: stmp1.Digest + "a"})
	if err != nil {
		t.Fatal(err)
	} else if key4 != "kore ga kagi dayo." || stmp4 == nil {
		t.Error(key4, stmp4)
	}
}
