package driver

import (
	"reflect"
	"testing"
)

func TestEventTree(t *testing.T) {
	tree := eventTree{}
	tree.add("/", Handler{&HandlerElement{Url: "a"}})
	tree.add("/a", Handler{&HandlerElement{Url: "b"}})
	tree.add("/a/a", Handler{&HandlerElement{Url: "c"}})
	tree.add("/a/b", Handler{&HandlerElement{Url: "d"}})
	tree.add("/b", Handler{&HandlerElement{Url: "f"}})
	tree.add("/b/a", Handler{&HandlerElement{Url: "f"}})
	tree.add("/b/b/a", Handler{&HandlerElement{Url: "g"}})

	if len(tree) != 8 {
		t.Error(len(tree), tree)
	}

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
	tree := eventTree{}
	tree.add("/", Handler{&HandlerElement{Url: "a"}})
	tree.add("/a", Handler{&HandlerElement{Url: "b"}})
	tree.add("/a/a", Handler{&HandlerElement{Url: "c"}})
	tree.add("/a/b", Handler{&HandlerElement{Url: "d"}})
	tree.add("/b", Handler{&HandlerElement{Url: "f"}})
	tree.add("/b/a", Handler{&HandlerElement{Url: "f"}})
	tree.add("/b/b/a", Handler{&HandlerElement{Url: "g"}})

	cont := tree.toContainer()
	if len(cont) != 7 {
		t.Error(len(cont), cont, tree)
	}
	tree2 := eventTree{}
	tree2.fromContainer(cont)
	if !reflect.DeepEqual(tree, tree2) {
		t.Error(tree, tree2)
	}
}
