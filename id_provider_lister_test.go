package driver

import (
	"reflect"
	"testing"
	"time"
)

// 事前に、UUID a_b-c、名前 ABC、URI https://localhost:1234 で登録しとく。

// 非キャッシュ用。
func testIdProviderLister(t *testing.T, reg IdProviderLister) {
	idps, err := reg.IdProviders()
	if err != nil {
		t.Fatal(err)
	} else if len(idps) == 0 {
		t.Error("No id providers.")
	} else if idps[0].Uuid != "a_b-c" ||
		idps[0].Name != "ABC" ||
		idps[0].LoginUri != "https://localhost:1234" {
		t.Error(idps[0])
	}
}

// キャッシュ用。
func testDatedIdProviderLister(t *testing.T, reg DatedIdProviderLister) {

	idps := []*IdProvider{
		&IdProvider{Uuid: "a_b-c", Name: "ABC", LoginUri: "https://localhost:1234"},
	}

	idps1, stmp1, err := reg.StampedIdProviders(nil)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(idps, idps1) || stmp1 == nil {
		t.Error(idps, idps1, stmp1)
	}

	// キャッシュと同じだから返らない。
	idps2, stmp2, err := reg.StampedIdProviders(stmp1)
	if err != nil {
		t.Fatal(err)
	} else if idps2 != nil || stmp2 == nil {
		t.Error(stmp1, idps2, stmp2)
	}

	// キャッシュが古いから返る。
	idps3, stmp3, err := reg.StampedIdProviders(&Stamp{Date: stmp1.Date.Add(-time.Second), Digest: stmp1.Digest})
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(idps, idps3) || stmp3 == nil {
		t.Error(idps, idps3, stmp3)
	}

	// ダイジェストが違うから返る。
	idps4, stmp4, err := reg.StampedIdProviders(&Stamp{Date: stmp1.Date, Digest: stmp1.Digest + "a"})
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(idps, idps4) || stmp4 == nil {
		t.Error(idps, idps4, stmp4)
	}
}
