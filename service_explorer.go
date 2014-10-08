package driver

import (
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"strings"
)

type ServiceExplorer interface {
	// サービスの URI から UUID を引く。
	ServiceUuid(servUri string, caStmp *Stamp) (servUuid string, newCaStmp *Stamp, err error)
}

// 骨組み。
type serviceExplorer struct {
	base KeyValueStore
}

func newServiceExplorer(base KeyValueStore) *serviceExplorer {
	return &serviceExplorer{base}
}

func (reg *serviceExplorer) ServiceUuid(servUri string, caStmp *Stamp) (servUuid string, newCaStmp *Stamp, err error) {
	value, newCaStmp, err := reg.base.Get("list", caStmp)
	if err != nil {
		return "", nil, erro.Wrap(err)
	} else if value == nil {
		return "", newCaStmp, nil
	}
	servUuid = value.(*serviceExplorerTree).get(servUri)
	if servUuid == "" {
		return "", nil, nil
	}
	return servUuid, newCaStmp, nil
}

// 内部データ。
type serviceExplorerTree struct {
	*util.Tree
}

func newServiceExplorerTree() *serviceExplorerTree {
	return &serviceExplorerTree{util.NewTree(serviceExplorerTreeIsRoot, serviceExplorerTreeParent)}
}

func serviceExplorerTreeIsRoot(label string) bool {
	return label == ""
}

func serviceExplorerTreeParent(label string) string {
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

func (tree *serviceExplorerTree) add(servUri string, servUuid string) {
	tree.Add(servUri, servUuid)
}

func (tree *serviceExplorerTree) remove(servUri string) {
	tree.Remove(servUri)
}

func (tree *serviceExplorerTree) get(servUri string) (servUuid string) {
	value := tree.ParentValue(servUri)
	if value == nil {
		return ""
	}
	return value.(string)
}

func (tree *serviceExplorerTree) fromContainer(cont map[string]string) {
	c := map[string]interface{}{}
	for name, addr := range cont {
		c[name] = addr
	}
	tree.FromContainer(c)
}

func (tree *serviceExplorerTree) toContainer() (cont map[string]string) {
	c := tree.ToContainer()
	cont = map[string]string{}
	for label, value := range c {
		cont[label] = value.(string)
	}
	return cont
}
