package driver

import (
	"reflect"
	"testing"
	"time"
)

func testJsBackendRegistry(t *testing.T, reg JsBackendRegistry) {
	dir := "/a/b"
	objName := "a_b"
	obj := &Object{true, true, []string{"$$http"}, "{a:function(){return 1+1}}"}

	obj1, stmp1, err := reg.StampedObject(dir, objName, nil)
	if err != nil {
		t.Fatal(err)
	} else if obj1 != nil || stmp1 != nil {
		t.Error(obj1, stmp1)
	}

	if err := reg.AddObject(dir, objName, obj); err != nil {
		t.Fatal(err)
	}

	// キャッシュの作成日時が対象の更新日時より後になるように待つ。
	timeUnit := time.Second // HTTP の If-Modified-Since とかを使っている場合、精度は秒。
	time.Sleep(timeUnit)

	obj2, stmp2, err := reg.StampedObject(dir, objName, nil)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(obj, obj2) || stmp2 == nil {
		t.Error(obj, obj2, stmp2)
	}

	// キャッシュと同じだから返らない。
	obj3, stmp3, err := reg.StampedObject(dir, objName, &Stamp{Date: stmp2.Date, Digest: stmp2.Digest})
	if err != nil {
		t.Fatal(err)
	} else if obj3 != nil || stmp3 == nil {
		t.Error(obj3, stmp3, stmp2)
	}

	// キャッシュが古いから返る。
	obj4, stmp4, err := reg.StampedObject(dir, objName, &Stamp{Date: stmp2.Date.Add(-2 * timeUnit), Digest: stmp2.Digest})
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(obj, obj4) || stmp4 == nil {
		t.Error(obj, obj4, stmp4, stmp2)
	}

	// ダイジェストが違うから返る。
	obj5, stmp5, err := reg.StampedObject(dir, objName, &Stamp{Date: stmp2.Date, Digest: stmp2.Digest + "a"})
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(obj, obj5) || stmp5 == nil {
		t.Error(obj, obj5, stmp5, stmp2)
	}

	if err = reg.RemoveObject(dir, objName); err != nil {
		t.Fatal(err)
	}

	obj6, stmp6, err := reg.StampedObject(dir, objName, nil)
	if err != nil {
		t.Fatal(err)
	} else if obj6 != nil || stmp6 != nil {
		t.Error(obj6, stmp6)
	}
}

// 事前に、UUID a_b-c、名前 ABC、URI https://localhost:1234 で登録しとく。
func testDatedIdProviderLister(t *testing.T, reg DatedIdProviderLister) {
	idps := []*IdProvider{
		&IdProvider{Uuid: "a_b-c", Name: "ABC", Uri: "https://localhost:1234"},
	}

	idps1, stmp1, err := reg.StampedIdProviders(nil)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(idps, idps1) || stmp1 == nil {
		t.Error(idps, idps1, stmp1)
	}

	// キャッシュと同じだから返らない。
	log.Debug("Aho")
	idps2, stmp2, err := reg.StampedIdProviders(stmp1)
	log.Debug("Baka")
	if err != nil {
		t.Fatal(err)
	} else if idps2 != nil || stmp2 == nil {
		t.Error(idps2, stmp2)
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
