package driver

import (
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"strings"
)

type TaExplorer interface {
	// サービスの URI から UUID を引く。
	ServiceUuid(servUri string, caStmp *Stamp) (servUuid string, newCaStmp *Stamp, err error)
}

// 骨組み。
type taExplorer struct {
	base KeyValueStore
}

func newTaExplorer(base KeyValueStore) *taExplorer {
	return &taExplorer{base}
}

func (reg *taExplorer) ServiceUuid(servUri string, caStmp *Stamp) (servUuid string, newCaStmp *Stamp, err error) {
	value, newCaStmp, err := reg.base.Get("list", caStmp)
	if err != nil {
		return "", nil, erro.Wrap(err)
	} else if value == nil {
		return "", newCaStmp, nil
	}
	servUuid = value.(*taExplorerTree).get(servUri)
	if servUuid == "" {
		return "", nil, nil
	}
	return servUuid, newCaStmp, nil
}

// 内部データ。
type taExplorerTree struct {
	*util.Tree
}

func newTaExplorerTree() *taExplorerTree {
	return &taExplorerTree{util.NewTree(taExplorerTreeIsRoot, taExplorerTreeParent)}
}

func taExplorerTreeIsRoot(label string) bool {
	return label == ""
}

func taExplorerTreeParent(label string) string {
	if idx := strings.LastIndex(label, "/"); idx < 0 {
		// localhost とか。
		return ""
	} else if sepIdx := strings.Index(label, "://"); sepIdx < 0 {
		if idx == len(label)-1 {
			// localhost/api/hoge/ とか。
			return label[:idx]
		} else {
			// localhost/api/hoge とか。
			return label[:idx+1]
		}
	} else if idx <= sepIdx+3 {
		// https:// とか
		return ""
	} else {
		// https://localhost/api/hoge とか。
		if idx == len(label)-1 {
			// localhost/api/hoge/ とか。
			return label[:idx]
		} else {
			// localhost/api/hoge とか。
			return label[:idx+1]
		}
	}
}

func (tree *taExplorerTree) add(servUri string, servUuid string) {
	tree.Add(servUri, servUuid)
}

func (tree *taExplorerTree) remove(servUri string) {
	tree.Remove(servUri)
}

func (tree *taExplorerTree) get(servUri string) (servUuid string) {
	value := tree.ParentValue(servUri)
	if value == nil {
		return ""
	}
	return value.(string)
}

func (tree *taExplorerTree) fromContainer(cont map[string]string) {
	c := map[string]interface{}{}
	for name, addr := range cont {
		c[name] = addr
	}
	tree.FromContainer(c)
}

func (tree *taExplorerTree) toContainer() (cont map[string]string) {
	c := tree.ToContainer()
	cont = map[string]string{}
	for label, value := range c {
		cont[label] = value.(string)
	}
	return cont
}
