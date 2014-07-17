package driver

import (
	"fmt"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net"
	"strings"
)

// NameRegistry の内部データ。

type nameTree map[string]*nameNode

type nameNode struct {
	addr       string
	childNames map[string]bool
}

func (node *nameNode) String() string {
	return fmt.Sprint(node.addr, toStringSlice(node.childNames))
}

func toStringSlice(set map[string]bool) []string {
	slice := []string{}
	for str, _ := range set {
		slice = append(slice, str)
	}
	return slice
}

func newNameNode() *nameNode {
	return &nameNode{childNames: map[string]bool{}}
}

// 通信形式から変換する。
func (tree nameTree) fromContainer(cont map[string]string) {
	for name, addr := range cont {
		tree.add(name, addr)
	}
}

// 通信形式へ変換する。
func (tree nameTree) toContainer() (cont map[string]string) {
	cont = map[string]string{}
	for name, node := range tree {
		if node.addr != "" {
			cont[name] = node.addr
		}
	}
	return cont
}

// 使用形式の木にノードを加える。
// 親が居なければ勝手に生成する。
func (tree nameTree) add(name string, addr string) {
	curNode := tree[name]
	if curNode == nil {
		curNode = newNameNode()
		tree[name] = curNode
	}

	curNode.addr = addr

	// 親を作成。
	for curName := name; curName != ""; {
		var parentName string
		if idx := strings.Index(curName, "."); idx < 0 {
			parentName = ""
		} else {
			parentName = curName[idx+1:]
		}
		parent := tree[parentName]
		if parent == nil {
			parent = newNameNode()
			tree[parentName] = parent
		}
		parent.childNames[curName] = true
		curName = parentName
	}
}

// 木からノードを削除する。
// ただし、子がいる場合は中間ノードとして残す。
func (tree nameTree) remove(name string) {
	node := tree[name]
	if node == nil {
		return
	}
	for childName, _ := range node.childNames {
		if tree[childName] != nil {
			node.addr = ""
			return
		}
	}
	delete(tree, name)

	// 親を削除。
	for curName := name; curName != ""; {
		var parentName string
		if idx := strings.Index(curName, "."); idx < 0 {
			parentName = ""
		} else {
			parentName = curName[idx+1:]
		}
		parent := tree[parentName]
		delete(parent.childNames, curName)
		if parent.addr != "" || len(parent.childNames) > 0 {
			break
		}
		curName = parentName
	}
}

func (tree nameTree) address(name string) (addr string) {
	node := tree[name]
	if node == nil {
		return ""
	}
	return node.addr
}

// 使用形式の木から部分木 + αに含まれるアドレスを列挙する。
func (tree nameTree) addresses(name string) (addrs []string) {
	node := tree[name]
	if node == nil {
		return nil
	}

	addrSet := map[string]bool{}
	tree._subTree(addrSet, name)
	return toStringSlice(addrSet)
}

func (tree nameTree) _subTree(addrSet map[string]bool, name string) {
	node := tree[name]
	if node == nil {
		return
	}

	if node.addr != "" {
		addrSet[node.addr] = true
	}
	for childName, _ := range node.childNames {
		tree._subTree(addrSet, childName)
	}
}

// 別名を登録アドレス（とポート）に展開する。
// 別名にポートが指定されている場合は、登録アドレスのポートを上書きする。
func ExpandName(reg NameRegistry, name string) (addrs []string, err error) {

	var portPart string
	if idx := strings.Index(name, ":"); idx >= 0 {
		portPart = name[idx:]
		name = name[:idx]
	}

	addrs, err = reg.Addresses(name)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	if portPart == "" {
		return addrs, nil
	}

	for i := 0; i < len(addrs); i++ {
		host, _, err := net.SplitHostPort(addrs[i])
		if err != nil {
			// 登録アドレスのポートは指定されてなかった。
			addrs[i] += portPart
		} else {
			addrs[i] = host + portPart
		}
	}

	return addrs, nil
}
