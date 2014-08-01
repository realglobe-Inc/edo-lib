package driver

import (
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net"
	"strings"
)

type NameRegistry interface {
	// アドレスを引く。
	Address(name string) (addr string, err error)
	// name はドメイン形式（. 区切りで後ろが親）の木構造のノードを表し、そのノード以下の部分木に含まれる全てのアドレスを返す。
	Addresses(name string) (addrs []string, err error)
}

// NameRegistry の内部データ。

type nameTree struct {
	*util.Tree
}

func newNameTree() *nameTree {
	return &nameTree{util.NewTree(
		func(label string) bool {
			return label == ""
		},
		func(label string) string {
			if idx := strings.Index(label, "."); idx < 0 {
				return ""
			} else {
				return label[idx+1:]
			}
		},
	)}
}

func (tree nameTree) add(name string, addr string) {
	tree.Add(name, addr)
}

func (tree nameTree) remove(name string) {
	tree.Remove(name)
}

func (tree nameTree) address(name string) (addr string) {
	val := tree.Value(name)
	if val == nil {
		return ""
	}
	return val.(string)
}

func (tree nameTree) addresses(name string) (addrs []string) {
	vals := tree.Values(name)
	if vals == nil {
		return nil
	}
	addrs = []string{}
	for _, val := range vals {
		addrs = append(addrs, val.(string))
	}
	return addrs
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

func (tree *nameTree) fromContainer(cont map[string]string) {
	c := map[string]interface{}{}
	for name, addr := range cont {
		c[name] = addr
	}
	tree.FromContainer(c)
}

func (tree *nameTree) toContainer() (cont map[string]string) {
	c := tree.ToContainer()
	cont = map[string]string{}
	for label, val := range c {
		cont[label] = val.(string)
	}
	return cont
}
