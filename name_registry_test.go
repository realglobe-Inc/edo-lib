package driver

import (
	"reflect"
	"testing"
)

func TestNameTree(t *testing.T) {
	tree := newNameTree()

	tree.add("", "a")
	tree.add("a", "b")
	tree.add("a.a", "c")
	tree.add("b.a", "d")
	tree.add("b", "e")
	tree.add("a.b", "f")
	tree.add("a.b.b", "g")

	if addrs := tree.addresses(""); len(addrs) != 7 {
		t.Error(addrs, tree)
	}
	if addrs := tree.addresses("a"); len(addrs) != 3 {
		t.Error(addrs, tree)
	}

	tree.remove("a")
	if addrs := tree.addresses("a"); len(addrs) != 2 {
		t.Error(addrs, tree)
	}
}

func TestNameTreeConversion(t *testing.T) {
	tree := newNameTree()
	tree.add("", "a")
	tree.add("a", "b")
	tree.add("a.a", "c")
	tree.add("b.a", "d")
	tree.add("b", "e")
	tree.add("a.b", "f")
	tree.add("a.b.b", "g")

	cont := tree.toContainer()
	if len(cont) != 7 {
		t.Error(len(cont), cont, tree)
	}
	tree2 := newNameTree()
	tree2.fromContainer(cont)

	if addrs := tree2.addresses(""); len(addrs) != 7 {
		t.Error(addrs, tree2)
	}
	if addrs := tree2.addresses("a"); len(addrs) != 3 {
		t.Error(addrs, tree2)
	}

	tree2.remove("a")
	if addrs := tree2.addresses("a"); len(addrs) != 2 {
		t.Error(addrs, tree2)
	}
}

var testNameTree *nameTree

func init() {
	testNameTree = newNameTree()
	testNameTree.fromContainer(map[string]string{
		"c.b.a": "c.localhost",
		"d.b.a": "d.localhost",
		"b.a":   "localhost",
	})
}

func testNameRegistry(t *testing.T, reg NameRegistry) {
	addr, _, err := reg.Address("c.b.a", nil)
	if err != nil {
		t.Fatal(err)
	} else if addr != "c.localhost" {
		t.Error(addr)
	}

	addrs, _, err := reg.Addresses("a", nil)
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
