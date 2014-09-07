package driver

import (
	"testing"
	"time"
)

// 事前に、UUID a_b-c、ユーザー属性取得用 URI https://localhost:1234/query で登録しとく。

// 非キャッシュ用。
func testIdProviderRegistry(t *testing.T, reg IdProviderRegistry) {
	queryUri, err := reg.IdProviderQueryUri("a_b-c")
	if err != nil {
		t.Fatal(err)
	} else if queryUri != "https://localhost:1234/query" {
		t.Error(queryUri)
	}
}

// キャッシュ用。
func testDatedIdProviderRegistry(t *testing.T, reg DatedIdProviderRegistry) {

	queryUri1, stmp1, err := reg.StampedIdProviderQueryUri("a_b-c", nil)
	if err != nil {
		t.Fatal(err)
	} else if queryUri1 != "https://localhost:1234/query" || stmp1 == nil {
		t.Error(queryUri1, stmp1)
	}

	// キャッシュと同じだから返らない。
	queryUri2, stmp2, err := reg.StampedIdProviderQueryUri("a_b-c", stmp1)
	if err != nil {
		t.Fatal(err)
	} else if queryUri2 != "" || stmp2 == nil {
		t.Error(queryUri2, stmp2)
	}

	// キャッシュが古いから返る。
	queryUri3, stmp3, err := reg.StampedIdProviderQueryUri("a_b-c", &Stamp{Date: stmp1.Date.Add(-time.Second), Digest: stmp1.Digest})
	if err != nil {
		t.Fatal(err)
	} else if queryUri3 != "https://localhost:1234/query" || stmp3 == nil {
		t.Error(queryUri3, stmp3)
	}

	// ダイジェストが違うから返る。
	queryUri4, stmp4, err := reg.StampedIdProviderQueryUri("a_b-c", &Stamp{Date: stmp1.Date, Digest: stmp1.Digest + "a"})
	if err != nil {
		t.Fatal(err)
	} else if queryUri4 != "https://localhost:1234/query" || stmp4 == nil {
		t.Error(queryUri4, stmp4)
	}
}
