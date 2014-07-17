package driver

import (
	"reflect"
	"testing"
	"time"
)

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

func testUserRegistry(t *testing.T, reg UserRegistry) {
	usrUuid := "a_b-c"
	attrName := "a b*c/d"
	attr := map[string]interface{}{"a": "b", "c": map[string]interface{}{"d": 1.08}}

	attrs1, err := reg.Attributes(usrUuid)
	if err != nil {
		t.Fatal(err)
	} else if len(attrs1) != 0 {
		t.Error(attrs1)
	}
	attr1, err := reg.Attribute(usrUuid, attrName)
	if err != nil {
		t.Fatal(err)
	} else if attr1 != nil {
		t.Error(attr1)
	}

	if err := reg.AddAttribute(usrUuid, attrName, attr); err != nil {
		t.Fatal(err)
	}

	attrs2, err := reg.Attributes(usrUuid)
	if err != nil {
		t.Fatal(err)
	} else if attrs2 == nil {
		t.Error(attrs2)
	}
	attr2, err := reg.Attribute(usrUuid, attrName)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(attr, attr2) {
		t.Error(attr2)
	}

	if err = reg.RemoveAttribute(usrUuid, attrName); err != nil {
		t.Fatal(err)
	}

	attrs3, err := reg.Attributes(usrUuid)
	if err != nil {
		t.Fatal(err)
	} else if len(attrs3) != 0 {
		t.Error(attrs3)
	}
	attr3, err := reg.Attribute(usrUuid, attrName)
	if err != nil {
		t.Fatal(err)
	} else if attr3 != nil {
		t.Error(attr3)
	}
}

func testJobRegistry(t *testing.T, reg JobRegistry) {
	usrUuid := "a_b-c"
	var jobId uint64 = 123
	deadline := time.Now().Add(time.Second)
	res := map[string]interface{}{"e": "f", "g": map[string]interface{}{"e": 1.08}}

	res1, err := reg.Result(usrUuid, jobId)
	if err != nil {
		t.Fatal(err)
	} else if res1 != nil {
		t.Error(res1)
	}

	if err := reg.AddResult(usrUuid, jobId, res, deadline); err != nil {
		t.Fatal(err)
	}

	res2, err := reg.Result(usrUuid, jobId)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(res, res2) {
		t.Error(res2)
	}
}

func testEventRegistry(t *testing.T, reg EventRegistry) {
	usrUuid := "a_b-c"
	event := "/d/e"
	var hndl Handler = []*HandlerElement{&HandlerElement{Url: "https://localhost"}}

	hndl1, err := reg.Handler(usrUuid, event)
	if err != nil {
		t.Fatal(err)
	} else if hndl1 != nil {
		t.Error(hndl1)
	}

	if err := reg.AddHandler(usrUuid, event, hndl); err != nil {
		t.Fatal(err)
	}

	hndl2, err := reg.Handler(usrUuid, event)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(hndl, hndl2) {
		t.Error(hndl, hndl2)
	}

	if err = reg.RemoveHandler(usrUuid, event); err != nil {
		t.Fatal(err)
	}

	hndl3, err := reg.Handler(usrUuid, event)
	if err != nil {
		t.Fatal(err)
	} else if hndl3 != nil {
		t.Error(hndl3)
	}
}
