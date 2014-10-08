package driver

import (
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net"
	"strings"
)

// 非キャッシュ用。
type NameRegistry interface {
	// アドレスを引く。
	Address(name string, caStmp *Stamp) (addr string, newCaStmp *Stamp, err error)
	// name はドメイン形式（. 区切りで後ろが親）の木構造のノードを表し、そのノード以下の部分木に含まれる全てのアドレスを返す。
	Addresses(name string, caStmp *Stamp) (addrs []string, newCaStmp *Stamp, err error)
}

// 骨組み。
type nameRegistry struct {
	base KeyValueStore
}

func newNameRegistry(base KeyValueStore) *nameRegistry {
	return &nameRegistry{base}
}

func (reg *nameRegistry) Address(name string, caStmp *Stamp) (addr string, newCaStmp *Stamp, err error) {
	value, newCaStmp, err := reg.base.Get("names", caStmp)
	if err != nil {
		return "", nil, erro.Wrap(err)
	} else if value == nil {
		return "", newCaStmp, nil
	}
	addr = value.(*nameTree).address(name)
	if addr == "" {
		return "", nil, nil
	}
	return addr, newCaStmp, nil
}

func (reg *nameRegistry) Addresses(name string, caStmp *Stamp) (addrs []string, newCaStmp *Stamp, err error) {
	value, newCaStmp, err := reg.base.Get("names", caStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if value == nil {
		return nil, newCaStmp, nil
	}
	addrs = value.(*nameTree).addresses(name)
	if addrs == nil {
		return nil, nil, nil
	}
	return addrs, newCaStmp, nil
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

func (tree *nameTree) add(name string, addr string) {
	tree.Add(name, addr)
}

func (tree *nameTree) remove(name string) {
	tree.Remove(name)
}

func (tree *nameTree) address(name string) (addr string) {
	value := tree.Value(name)
	if value == nil {
		return ""
	}
	return value.(string)
}

func (tree *nameTree) addresses(name string) (addrs []string) {
	values := tree.Values(name)
	if values == nil {
		return nil
	}
	addrs = []string{}
	for _, value := range values {
		addrs = append(addrs, value.(string))
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

	addrs, _, err = reg.Addresses(name, nil)
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
	for label, value := range c {
		cont[label] = value.(string)
	}
	return cont
}
