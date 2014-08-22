package driver

import (
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
