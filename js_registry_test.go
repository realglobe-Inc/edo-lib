package driver

import (
	"reflect"
	"testing"
	"time"
)

// 非キャッシュ用。
func testJsRegistry(t *testing.T, reg JsRegistry) {
	dir := "/a/b"
	objName := "a_b"
	obj := &Object{true, true, []string{"$$http"}, "{a:function(){return 1+1}}"}

	obj1, err := reg.Object(dir, objName)
	if err != nil {
		t.Fatal(err)
	} else if obj1 != nil {
		t.Error(obj1)
	}

	if err := reg.AddObject(dir, objName, obj); err != nil {
		t.Fatal(err)
	}

	obj2, err := reg.Object(dir, objName)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(obj, obj2) {
		t.Error(obj, obj2)
	}

	if err = reg.RemoveObject(dir, objName); err != nil {
		t.Fatal(err)
	}

	obj3, err := reg.Object(dir, objName)
	if err != nil {
		t.Fatal(err)
	} else if obj3 != nil {
		t.Error(obj3)
	}
}

// キャッシュ用。
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
