package driver

import (
	"github.com/realglobe-Inc/edo/util"
	"strconv"
	"strings"
	"time"
)

// メモリ上で完結する。デバッグ用。

// 非キャッシュ用。
type MemoryServiceExplorer struct {
	*serviceExplorerTree
}

func NewMemoryServiceExplorer() *MemoryServiceExplorer {
	return &MemoryServiceExplorer{newServiceExplorerTree()}
}

func (reg *MemoryServiceExplorer) ServiceUuid(servUri string) (servUuid string, err error) {
	return reg.get(servUri), nil
}
func (reg *MemoryServiceExplorer) AddServiceUuid(servUri, servUuid string) {
	reg.add(servUri, servUuid)
}
func (reg *MemoryServiceExplorer) RemoveIdProvider(servUri string) {
	reg.remove(servUri)
}

// キャッシュ用。
type MemoryDatedServiceExplorer struct {
	*MemoryServiceExplorer
	stmp    *Stamp
	expiDur time.Duration
}

func NewMemoryDatedServiceExplorer(expiDur time.Duration) *MemoryDatedServiceExplorer {
	return &MemoryDatedServiceExplorer{NewMemoryServiceExplorer(), &Stamp{Date: time.Now(), Digest: strconv.Itoa(0)}, expiDur}
}

func (reg *MemoryDatedServiceExplorer) StampedServiceUuid(servUri string, caStmp *Stamp) (servUuid string, newCaStmp *Stamp, err error) {
	newCaStmp = &Stamp{Date: reg.stmp.Date, ExpiDate: time.Now().Add(reg.expiDur), Digest: reg.stmp.Digest}

	if caStmp == nil || caStmp.Date.Before(reg.stmp.Date) || caStmp.Digest != reg.stmp.Digest {
		servUuid, _ = reg.ServiceUuid(servUri)
		if servUuid == "" {
			return "", nil, nil
		} else {
			return servUuid, newCaStmp, nil
		}
	}

	return "", newCaStmp, nil
}
func (reg *MemoryDatedServiceExplorer) AddServiceUuid(servUri, servUuid string) {
	reg.MemoryServiceExplorer.AddServiceUuid(servUri, servUuid)
	dig, _ := strconv.Atoi(reg.stmp.Digest)
	reg.stmp = &Stamp{Date: time.Now(), Digest: strconv.Itoa(dig + 1)}
}
func (reg *MemoryDatedServiceExplorer) RemoveServiceUuid(servUri string) {
	reg.MemoryServiceExplorer.RemoveIdProvider(servUri)
	dig, _ := strconv.Atoi(reg.stmp.Digest)
	reg.stmp = &Stamp{Date: time.Now(), Digest: strconv.Itoa(dig + 1)}
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
			// localhost/api/hoge/ とか。
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
			// localhost/api/hoge/ とか。
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
	val := tree.ParentValue(servUri)
	if val == nil {
		return ""
	}
	return val.(string)
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
	for label, val := range c {
		cont[label] = val.(string)
	}
	return cont
}
