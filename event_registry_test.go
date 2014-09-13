package driver

import (
	"reflect"
	"testing"
)

func TestEventTree(t *testing.T) {
	tree := newEventTree()
	tree.add("/", Handler{&HandlerElement{Url: "a"}})
	tree.add("/a", Handler{&HandlerElement{Url: "b"}})
	tree.add("/a/a", Handler{&HandlerElement{Url: "c"}})
	tree.add("/a/b", Handler{&HandlerElement{Url: "d"}})
	tree.add("/b", Handler{&HandlerElement{Url: "f"}})
	tree.add("/b/a", Handler{&HandlerElement{Url: "g"}})
	tree.add("/b/b/a", Handler{&HandlerElement{Url: "h"}})

	if hndl := tree.handler("/"); len(hndl) != 7 {
		t.Error(hndl, tree)
	}
	if hndl := tree.handler("/a"); len(hndl) != 3 {
		t.Error(hndl, tree)
	}

	tree.remove("/a")
	if hndl := tree.handler("/a"); len(hndl) != 2 {
		t.Error(hndl, tree)
	}
}

func TestEventTreeConversion(t *testing.T) {
	tree := newEventTree()
	tree.add("/", Handler{&HandlerElement{Url: "a"}})
	tree.add("/a", Handler{&HandlerElement{Url: "b"}})
	tree.add("/a/a", Handler{&HandlerElement{Url: "c"}})
	tree.add("/a/b", Handler{&HandlerElement{Url: "d"}})
	tree.add("/b", Handler{&HandlerElement{Url: "f"}})
	tree.add("/b/a", Handler{&HandlerElement{Url: "g"}})
	tree.add("/b/b/a", Handler{&HandlerElement{Url: "h"}})

	cont := tree.toContainer()
	if len(cont) != 7 {
		t.Error(len(cont), cont, tree)
	}
	tree2 := newEventTree()
	tree2.fromContainer(cont)

	if hndl := tree2.handler("/"); len(hndl) != 7 {
		t.Error(hndl, tree2)
	}
	if hndl := tree2.handler("/a"); len(hndl) != 3 {
		t.Error(hndl, tree2)
	}

	tree2.remove("/a")
	if hndl := tree2.handler("/a"); len(hndl) != 2 {
		t.Error(hndl, tree2)
	}
}

// 非キャッシュ用。
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
