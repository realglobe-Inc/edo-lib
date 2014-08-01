package driver

import (
	"reflect"
	"testing"
	"time"
)

// 事前に、
// abc-012 に a_b-c、
// を登録しとく。
func testLoginRegistry(t *testing.T, reg LoginRegistry) {
	usrUuid, err := reg.User("abc-012")
	if err != nil {
		t.Fatal(err)
	} else if usrUuid != "a_b-c" {
		t.Error(usrUuid)
	}
}

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
		t.Error(attr, attr2)
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
	var jobId string = "a_b-c"
	deadline := time.Now().Add(time.Second)
	res := &JobResult{Status: 200, Body: "konnnann demashita"}

	res1, err := reg.Result(jobId)
	if err != nil {
		t.Fatal(err)
	} else if res1 != nil {
		t.Error(res1)
	}

	if err := reg.AddResult(jobId, res, deadline); err != nil {
		t.Fatal(err)
	}

	res2, err := reg.Result(jobId)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(res, res2) {
		t.Error(res2)
	}
}

func testJobRegistryRemoveOld(t *testing.T, reg JobRegistry) {
	jobId := "a_b-c"
	deadline := time.Now().Add(50 * time.Millisecond)
	res := &JobResult{Status: 200, Body: "konnnann demashita"}

	if err := reg.AddResult(jobId+"1", res, deadline); err != nil {
		t.Fatal(err)
	}
	if err := reg.AddResult(jobId+"2", res, deadline.Add(time.Second)); err != nil {
		t.Fatal(err)
	}

	time.Sleep(100 * time.Millisecond)

	if err := reg.AddResult(jobId+"3", res, deadline); err != nil {
		t.Fatal(err)
	}

	res1, err := reg.Result(jobId + "1")
	if err != nil {
		t.Fatal(err)
	} else if res1 != nil {
		t.Error(res1)
	}
	res2, err := reg.Result(jobId + "2")
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(res, res2) {
		t.Error(res, res2)
	}

}

// 事前に、
// c.b.a に c.localhost、
// d.b.a に d.localhost、
// b.a   に   localhost、
// を登録しとく。
func testNameRegistry(t *testing.T, reg NameRegistry) {
	addr, err := reg.Address("c.b.a")
	if err != nil {
		t.Fatal(err)
	} else if addr != "c.localhost" {
		t.Error(addr)
	}

	addrs, err := reg.Addresses("a")
	if err != nil {
		t.Fatal(err)
	}
	set := map[string]bool{}
	for _, addr := range addrs {
		set[addr] = true
	}
	if !reflect.DeepEqual(map[string]bool{"c.localhost": true, "d.localhost": true, "localhost": true}, set) {
		t.Error(addrs)
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

// 事前に、
// localhost:1234 に a_b-c、
// を登録しとく。
func testServiceRegistry(t *testing.T, reg ServiceRegistry) {
	servUuid, err := reg.Service("localhost:1234")
	if err != nil {
		t.Fatal(err)
	} else if servUuid != "a_b-c" {
		t.Error(servUuid)
	}
}
