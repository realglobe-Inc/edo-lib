package driver

import (
	"testing"
	"time"
)

// 事前に、UUID a_b-c、ユーザー属性取得用 URI https://localhost:1234/query で登録しとく。

// 非キャッシュ用。
func testServiceExplorer(t *testing.T, reg ServiceExplorer) {
	servUuid, err := reg.ServiceUuid("https://localhost:1234/api/hoge")
	if err != nil {
		t.Fatal(err)
	} else if servUuid != "a_b-c" {
		t.Error(servUuid)
	}

	servUuid, err = reg.ServiceUuid("https://localhost:1234/api")
	if err != nil {
		t.Fatal(err)
	} else if servUuid != "a_b-c" {
		t.Error(servUuid)
	}

	servUuid, err = reg.ServiceUuid("https://localhost:1234")
	if err != nil {
		t.Fatal(err)
	} else if servUuid != "" {
		t.Error(servUuid)
	}
}

// キャッシュ用。
func testDatedServiceExplorer(t *testing.T, reg DatedServiceExplorer) {

	servUuid1, stmp1, err := reg.StampedServiceUuid("https://localhost:1234/api", nil)
	if err != nil {
		t.Fatal(err)
	} else if servUuid1 != "a_b-c" || stmp1 == nil {
		t.Error(servUuid1, stmp1)
	}

	// キャッシュと同じだから返らない。
	servUuid2, stmp2, err := reg.StampedServiceUuid("https://localhost:1234/api", stmp1)
	if err != nil {
		t.Fatal(err)
	} else if servUuid2 != "" || stmp2 == nil {
		t.Error(servUuid2, stmp2)
	}

	// キャッシュが古いから返る。
	servUuid3, stmp3, err := reg.StampedServiceUuid("https://localhost:1234/api", &Stamp{Date: stmp1.Date.Add(-time.Second), Digest: stmp1.Digest})
	if err != nil {
		t.Fatal(err)
	} else if servUuid3 != "a_b-c" || stmp3 == nil {
		t.Error(servUuid3, stmp3)
	}

	// ダイジェストが違うから返る。
	servUuid4, stmp4, err := reg.StampedServiceUuid("https://localhost:1234/api", &Stamp{Date: stmp1.Date, Digest: stmp1.Digest + "a"})
	if err != nil {
		t.Fatal(err)
	} else if servUuid4 != "a_b-c" || stmp4 == nil {
		t.Error(servUuid4, stmp4)
	}
}
