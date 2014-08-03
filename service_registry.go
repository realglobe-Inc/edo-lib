package driver

import (
	"github.com/realglobe-Inc/edo/util"
	"strings"
)

type ServiceRegistry interface {
	// EDO に登録されたサービスの管轄外向けエンドポイントから UUID を引く。
	Service(endPt string) (servUuid string, err error)
}

// ServiceRegistry の内部データ。

type serviceTree struct {
	*util.Tree
}

func newServiceTree() *serviceTree {
	return &serviceTree{util.NewTree(
		func(label string) bool {
			return label == ""
		},
		func(label string) string {
			if idx := strings.LastIndex(label, "/"); idx < 0 {
				// localhost とか。
				return ""
			} else if sepIdx := strings.Index(label, "://"); sepIdx < 0 {
				// localhost/api/hoge とか。
				return label[:idx]
			} else if idx <= sepIdx+3 {
				// https:// とか
				return ""
			} else {
				// https://localhost/api/hoge とか。
				return label[:idx]
			}
		},
	)}
}

func (tree *serviceTree) add(endPt string, servUuid string) {
	tree.Add(endPt, servUuid)
}

func (tree *serviceTree) remove(endPt string) {
	tree.Remove(endPt)
}

func (tree *serviceTree) service(endPt string) (servUuid string) {
	val := tree.ParentValue(endPt)
	if val == nil {
		return ""
	}
	return val.(string)
}

func (tree *serviceTree) fromContainer(cont map[string]string) {
	c := map[string]interface{}{}
	for name, addr := range cont {
		c[name] = addr
	}
	tree.FromContainer(c)
}

func (tree *serviceTree) toContainer() (cont map[string]string) {
	c := tree.ToContainer()
	cont = map[string]string{}
	for label, val := range c {
		cont[label] = val.(string)
	}
	return cont
}
